from django.db import models
from django.forms.models import model_to_dict
import json

class Inspiration_Table(models.Model):
    name = models.CharField(max_length=100)
    count = models.ManyToManyField('Roll', related_name='count')
    weighted_inspirations = models.ManyToManyField('Weighted_Inspiration')
    tags = models.ManyToManyField('Tag')

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
    tags = models.ManyToManyField('Tag')
    inspiration_tables = models.ManyToManyField(Inspiration_Table, related_name='inspiration_tables')
