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

    name = models.CharField(default="-", max_length=100)
    tags = models.ManyToManyField('Tag')
    flags = models.CharField(blank=True, null=True, max_length=100)
    muted_flags = models.CharField(blank=True, null=True, max_length=100)
    required_flags = models.CharField(blank=True, null=True, max_length=200)
    system = models.ForeignKey('Config_System', null=True, blank=True, on_delete=models.SET_NULL)
    star_cluster = models.ForeignKey('Config_Star_Cluster', null=True, blank=True, on_delete=models.SET_NULL)
    route = models.ForeignKey('Config_Route', null=True, blank=True, on_delete=models.SET_NULL)
    zone = models.ForeignKey('Config_Zone', null=True, blank=True, on_delete=models.SET_NULL, related_name='perterbation_zone')
    element = models.ForeignKey('Config_Element', null=True, blank=True, on_delete=models.SET_NULL)
    satellite = models.ForeignKey('Config_Element', null=True, blank=True, on_delete=models.SET_NULL, related_name='perterbation_satellite')
    territory = models.ForeignKey('Config_Territory', null=True, blank=True, on_delete=models.SET_NULL)
