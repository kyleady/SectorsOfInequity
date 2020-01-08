from django.db import models
from django.forms.models import model_to_dict
import json

from .perterbation import Perterbation_System

class BaseInspiration(models.Model):
    class Meta:
        abstract = True

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(default="-", max_length=25)
    description = models.CharField(default="-", max_length=1000)
    perterbation = None


class Inspiration_System_Feature(BaseInspiration):
    perterbation = models.ForeignKey(Perterbation_System, on_delete=models.CASCADE)
