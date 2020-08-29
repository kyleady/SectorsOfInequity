from django.db import models
from django.forms.models import model_to_dict
import json

class Inspiration_Extra(models.Model):
    name = models.CharField(max_length=100)
    count = models.ManyToManyField('Roll', related_name='extra_tables_count')
    type = models.ForeignKey('Config_Name', on_delete=models.CASCADE, related_name='inspiration_extra_type')
    inspiration_tables = models.ManyToManyField('Inspiration_Table', related_name='inspiration_tables')
    tags = models.ManyToManyField('Tag', related_name='inspiration_extra_tags')

class Inspiration_Table(models.Model):
    name = models.CharField(max_length=100)
    count = models.ManyToManyField('Roll', related_name='inspiration_table_count')
    modifiers = models.ManyToManyField('Roll', related_name='inspiration_table_modifiers')
    weighted_inspirations = models.ManyToManyField('Weighted_Inspiration', related_name='inspiration_table_weighted_inspirations')
    extra_inspirations = models.ManyToManyField('Weighted_Inspiration', related_name='inspiration_table_extra_inspirations')
    tags = models.ManyToManyField('Tag', related_name='inspiration_table_tags')

class Inspiration(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(max_length=75)
    description = models.TextField()
    roll_groups = models.ManyToManyField(Inspiration_Table, related_name='roll_groups')
    perterbations = models.ManyToManyField('Perterbation')
    tags = models.ManyToManyField('Tag', related_name='inspiration_tags')
    inspiration_tables = models.ManyToManyField(Inspiration_Table, related_name='sub_inspiration_tables')
