from django.db import models
from django.forms.models import model_to_dict
import json

# Abstract Models
class BaseConfig(models.Model):
    class Meta:
        abstract = True

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(default="-", max_length=100)

class Config_Territory(BaseConfig):
    inspirations = models.ManyToManyField('Inspiration_Table', related_name='territory_inspiration')

class Config_Element(BaseConfig):
    inspirations = models.ManyToManyField('Inspiration_Table', related_name='element_inspiration')
    satellite_count =  models.ManyToManyField('Roll', related_name='element_satellite_count')
    satellite_extra =  models.ManyToManyField('Extra_Tables', related_name='element_satellite_extra')
    territory_count =  models.ManyToManyField('Roll', related_name='element_territory_count')
    territory_extra =  models.ManyToManyField('Extra_Tables', related_name='element_territory_extra')
    spacing = models.ManyToManyField('Roll', related_name='element_spacing')
    territory = models.ForeignKey('Config_Territory', null=True, blank=True, on_delete=models.SET_NULL, related_name='element_territory')

class Config_Zone(BaseConfig):
    inspirations = models.ManyToManyField('Inspiration_Table', related_name='zone_inspiration')
    element_count =  models.ManyToManyField('Roll', related_name='zone_element_count')
    element_extra =  models.ManyToManyField('Extra_Tables', related_name='zone_element_extra')
    order = models.ManyToManyField('Roll', related_name='zone_order')

class Config_Star_Cluster(BaseConfig):
    inspirations = models.ManyToManyField('Inspiration_Table', related_name='star_cluster_inspiration')
    zone_count =  models.ManyToManyField('Roll', related_name='star_cluster_zone_count')
    zone_extra =  models.ManyToManyField('Extra_Tables', related_name='star_cluster_zone_extra')

class Config_Route(BaseConfig):
    inspirations = models.ManyToManyField('Inspiration_Table', related_name='route_inspiration')

class Config_System(BaseConfig):
    inspirations = models.ManyToManyField('Inspiration_Table', related_name='system_inspiration')
    star_cluster_count = models.ManyToManyField('Roll', related_name='star_cluster_count')
    star_cluster_extra = models.ManyToManyField('Extra_Tables', related_name='star_cluster_extra')

class Config_Grid(BaseConfig):
    height = models.PositiveSmallIntegerField(default=20, blank=True)
    width = models.PositiveSmallIntegerField(default=20, blank=True)
    connectionRange = models.PositiveSmallIntegerField(default=5, blank=True)
    populationRate = models.FloatField(default=0.5, blank=True)
    connectionRate = models.FloatField(default=0.5, blank=True)
    rangeRateMultiplier = models.FloatField(default=0.5, blank=True)
    smoothingFactor = models.FloatField(default=0.5, blank=True)
