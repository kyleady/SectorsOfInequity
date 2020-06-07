from django.db import models
from django.forms.models import model_to_dict
import json

from .perterbation import Perterbation
from .roll import Roll
from .tag import Tag

class Inspiration_Nested(models.Model):
    count = models.ManyToManyField(Roll, related_name='count')

class Inspiration(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(default="-", max_length=25)
    description = models.CharField(default="-", max_length=1000)
    rolls = models.ManyToManyField(Roll, related_name='rolls')
    perterbation = models.ForeignKey(Perterbation, null=True, blank=True, on_delete=models.CASCADE)
    tags = models.ManyToManyField(Tag)
    nested_inspirations = models.ManyToManyField(Inspiration_Nested, related_name='nested_inspirations')
