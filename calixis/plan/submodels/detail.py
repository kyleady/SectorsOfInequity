from django.db import models
from django.forms.models import model_to_dict
from django.core.validators import int_list_validator
import json
import re

class Detail(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.inspirations.all()[0].name

    def get_rolls(self):
        return self.rolls.split(',')

    def get_inspiration_table_name(self):
        inspiration_tables = self.inspiration_tables.all()
        if len(self.inspiration_tables.all()) > 0:
            return inspiration_tables[0].name
        else:
            return "-"

    def get_description(self):
        rolls = json.loads(self.rolls)
        text = '\n'.join(inspiration.description for inspiration in self.inspirations.all()).format(**rolls)
        return text

    def get_all_child_details(self):
        all_child_details = []
        for child_detail in self.detail_set.all():
            all_child_details.append(child_detail)
            all_child_details.extend(child_detail.get_all_child_details())

        return all_child_details

    rolls = models.CharField(max_length=100)
    inspirations = models.ManyToManyField('Inspiration', related_name='inspirations')
    inspiration_tables = models.ManyToManyField('Inspiration_Table')
    parent_detail = models.ForeignKey('self', null=True, blank=True, on_delete=models.CASCADE)
