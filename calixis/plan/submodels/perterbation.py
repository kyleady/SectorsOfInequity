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

    name = models.CharField(default="-", max_length=25)
    tags = models.ManyToManyField('Tag')
    system = models.ForeignKey('Config_System', null=True, blank=True, on_delete=models.SET_NULL)
    star_cluster = models.ForeignKey('Config_Star_Cluster', null=True, blank=True, on_delete=models.SET_NULL)
    route = models.ForeignKey('Config_Route', null=True, blank=True, on_delete=models.SET_NULL)
    zones = models.ManyToManyField('Config_Zone', related_name='zones')
    element = models.ForeignKey('Config_Element', null=True, blank=True, on_delete=models.SET_NULL)
