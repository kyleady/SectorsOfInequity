from django.db import models
from django.forms.models import model_to_dict
import json

class Inspiration_Nested(models.Model):
    name = models.CharField(max_length=25)
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

    name = models.CharField(max_length=25)
    description = models.CharField(default="", max_length=1000)
    roll_groups = models.ManyToManyField(Inspiration_Nested, related_name='roll_groups')
    perterbations = models.ManyToManyField('Perterbation')
    tags = models.ManyToManyField('Tag')
    nested_inspirations = models.ManyToManyField(Inspiration_Nested, related_name='nested_inspirations')
