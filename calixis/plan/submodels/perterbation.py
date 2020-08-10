from django.db import models
from django.forms.models import model_to_dict
import json

class Perterbation(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(default="", max_length=100)
    tags = models.ManyToManyField('Tag', related_name='perterbation_tags')
    flags = models.CharField(blank=True, null=True, max_length=100)
    muted_flags = models.CharField(blank=True, null=True, max_length=100)
    required_flags = models.CharField(blank=True, null=True, max_length=200)
    configs = models.ManyToManyField('Config_Asset', related_name='perterbation_configs')
