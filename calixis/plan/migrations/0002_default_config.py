from django.db import migrations, models

import json
import os
import pprint
import uuid
import yaml
import sys

pp = pprint.PrettyPrinter(indent=4)

def _fetch_instructions(directory):
    instructions = []
    path_to_script = os.path.abspath(os.path.dirname(__file__))
    path = os.path.join(path_to_script, directory)
    for filename in os.listdir(path):
        if filename.endswith(".yml") or filename.endswith(".yaml"):
            filename = os.path.join(path, filename)
            file = open(filename)
            yaml_instructions = yaml.load(file)
            instructions += yaml_instructions

    return instructions

def _get_plan(instructions):
    plan = {
        'references': [],
        'rows': [],
        'plan_ids': []
    }

    for instruction in instructions:
        if isinstance(instruction, dict):
            is_reference = False
            counts = {
                'reference': 0,
                'row': 0,
                'fk': 0,
                'm2m': 0
            }
            fk_dependencies = []
            m2m_dependencies = []

            if '_type' not in instruction:
                raise Exception('Model does not have a _type key')

            for k, v in instruction.items():
                sub_plan = {
                    'references': [],
                    'rows': [],
                    'top_level': []
                }
                if isinstance(v, dict):
                    sub_plan = _get_plan([instruction[k]])
                    instruction[k] = {'_plan_id': sub_plan['plan_ids'][0]}
                    fk_dependencies += sub_plan['plan_ids']
                    counts['fk'] += 1
                elif isinstance(v, list):
                    sub_plan = _get_plan(instruction[k])
                    instruction[k] = [{'_plan_id': plan_id} for plan_id in sub_plan['plan_ids']]
                    m2m_dependencies += sub_plan['plan_ids']
                    counts['m2m'] += len(v)
                elif k == '_exists' and v == True:
                    is_reference = True

                plan['references'] += sub_plan['references']
                plan['rows'] += sub_plan['rows']
                counts['reference'] += len(sub_plan['references'])
                counts['row'] += len(sub_plan['rows'])

            instruction['_fk_dependencies'] = fk_dependencies
            instruction['_m2m_dependencies'] = m2m_dependencies
            instruction['_plan_id'] = str(uuid.uuid4())
            plan['plan_ids'].append(instruction['_plan_id'])
            for count_type, count_total in counts.items():
                instruction['_{}_count'.format(count_type)] = count_total

            row_plan = {
                'args': { key:value for key, value in instruction.items() if not key.startswith('_') and not isinstance(value, list) },
                'metadata': { key:value for key, value in instruction.items() if key.startswith('_') },
                'm2m_args': { key:value for key, value in instruction.items() if not key.startswith('_') and isinstance(value, list) }
            }

            if is_reference:
                if counts['row'] > 0 or counts['reference'] > 0:
                    raise Exception('A reference cannot create other rows or reference other rows. Query the related rows through this reference instead.')
                plan['references'].append(row_plan)
            else:
                plan['rows'].append(row_plan)
        else:
            raise Exception('Model of type {} is not a dict.'.format(type(instruction)))

    return plan

def _dependencies_exist(row_plan, created_rows, dependency_type):
    all_dependencies_have_been_created = True
    for dependency in row_plan['metadata']['_{}_dependencies'.format(dependency_type)]:
        if dependency not in created_rows:
            all_dependencies_have_been_created = False
            break

    return all_dependencies_have_been_created

def _attempt_to_create_each(todo, created_rows, counts, model_constructors):
    for _plan_id, row_plan in list(todo['rows'].items()):
        if not _dependencies_exist(row_plan, created_rows, 'fk'):
            continue

        for arg, value in row_plan['args'].items():
            if isinstance(value, dict):
                row_plan['args'][arg] = created_rows[value['_plan_id']]
        try:
            created_row = model_constructors[row_plan['metadata']['_type']](**row_plan['args'])
            created_row.save()
        except:
            print(row_plan)
            print("Unexpected error:", sys.exc_info()[0])
            raise
        if len(row_plan['metadata']['_m2m_dependencies']) > 0:
            todo['m2m'][_plan_id] = (row_plan, created_row)

        created_rows[_plan_id] = created_row
        counts['rows'] += 1
        del todo['rows'][_plan_id]
        counts['activity'] += 1

        counts['types'][row_plan['metadata']['_type']] = counts['types'].get(row_plan['metadata']['_type'], 0) + 1

def _attempt_to_connect_each(todo, created_rows, counts):
    for _plan_id, (row_plan, created_row) in list(todo['m2m'].items()):
        if not _dependencies_exist(row_plan, created_rows, 'm2m'):
            continue

        for arg, ids in row_plan['m2m_args'].items():
            getattr(created_row, arg).add(*[created_rows[id['_plan_id']] for id in ids])

        del todo['m2m'][_plan_id]
        counts['activity'] += 1

def _attempt_to_reference_each(todo, created_rows, counts, model_constructors):
    for _plan_id, reference_plan in list(todo['references'].items()):
        try:
            found_objects = model_constructors[reference_plan['metadata']['_type']].objects.filter(**reference_plan['args'])
        except:
            pp.pprint(reference_plan['metadata']['_type'])
            pp.pprint(reference_plan['args'])
            raise Exception('Error occured while searching for a row.')

        number_of_found_objects = len(found_objects)
        if number_of_found_objects > 1 and not reference_plan['metadata'].get('_limit_1', False):
            pp.pprint(reference_plan['metadata']['_type'])
            pp.pprint(reference_plan['args'])
            raise Exception('Referenced row returned more than one result. Use the _limit_1 flag if you expect more than one result. Migration partially executed.')
        elif number_of_found_objects < 1:
            continue

        created_rows[_plan_id] = found_objects[0]
        del todo['references'][_plan_id]
        counts['activity'] += 1
        counts['references'] += 1

def _create_rows(apps, plan, default_app_name):
    model_constructors = { type:None for type in (
                [row_plan['metadata']['_type'] for row_plan in plan['rows']]
                 + [reference_plan['metadata']['_type'] for reference_plan in plan['references']]) }
    for _type in model_constructors:
        app_name = default_app_name
        model_name = _type
        if '.' in _type:
            (app_name, _name) = _type.split('.')
        model_constructors[_type] = apps.get_model(app_name, model_name)

    counts = {
        'rows': 0,
        'references': 0,
        'activity': -1,
        'types': {}
    }

    todo = {
        'rows': { plan_id:row_plan for plan_id, row_plan in [(row_plan['metadata']['_plan_id'], row_plan) for row_plan in plan['rows']] },
        'references': { plan_id:reference_plan for plan_id, reference_plan in [(reference_plan['metadata']['_plan_id'], reference_plan) for reference_plan in plan['references']] },
        'm2m': {}
    }

    created_rows = {}
    while counts['activity'] != 0:
        counts['activity'] = 0
        _attempt_to_create_each(todo, created_rows, counts, model_constructors)
        _attempt_to_connect_each(todo, created_rows, counts)
        _attempt_to_reference_each(todo, created_rows, counts, model_constructors)

    for todo_list in todo.values():
        if len(todo_list) > 0:
            print('REMINAING ROWS TO CREATE')
            print("Rows created: {}".format(counts['rows']))
            pp.pprint(counts['types'])
            print("Rows to create: {}".format(len(todo['rows'])))
            pp.pprint(todo['rows'])

            print('REMINAING ROWS TO REFERENCE')
            print("Rows referenced: {}".format(counts['references']))
            print("Rows to reference: {}".format(len(todo['references'])))
            pp.pprint(todo['references'])

            raise Exception('Instructions had impossible foreign key dependencies and/or references to rows that do not exist. Migration partially executed.')

def migrate_by_yml(apps, directory, default_app_name):
    instructions = _fetch_instructions(directory)
    plan = _get_plan(instructions)
    _create_rows(apps, plan, default_app_name)

def add_example_config(apps, schema_editor):
    migrate_by_yml(apps, '0002_default_config/', 'plan')

class Migration(migrations.Migration):
    dependencies = [
        ('plan', '0001_initial'),
    ]

    operations = [
        migrations.RunPython(add_example_config)
    ]

if __name__ == "__main__":
    directory = '0002_default_config/'
    instructions = _fetch_instructions(directory)
    plan = _get_plan(instructions)
    with open('debug.json', 'w') as outfile:
        json.dump(plan, outfile, sort_keys=True, indent=4, separators=(',', ': '))
        print('Plan written to debug.json.')
