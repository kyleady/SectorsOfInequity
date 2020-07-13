from django.db import models
from django.forms.models import model_to_dict
import json

# Abstract Models
class BaseAsset(models.Model):
    class Meta:
        abstract = True

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(max_length=75)
    parent = None

class Asset_Territory(BaseAsset):
    type = models.ForeignKey('Detail', on_delete=models.CASCADE, related_name='territory_type')

class Asset_Element(BaseAsset):
    type = models.ForeignKey('Detail', on_delete=models.CASCADE, related_name='element_type')
    distance = models.IntegerField()
    satellites = models.ManyToManyField('Asset_Element', related_name='element_satellites')
    territories = models.ManyToManyField('Asset_Territory', related_name='element_territories')

class Asset_Zone(BaseAsset):
    distance = models.SmallIntegerField()
    elements = models.ManyToManyField(Asset_Element, related_name='elements')

class Asset_Star_Cluster(BaseAsset):
    stars = models.ManyToManyField('Detail', related_name='stars')
    zones = models.ManyToManyField(Asset_Zone, related_name='zones')

class Asset_Route(BaseAsset):
    stability = models.ForeignKey('Detail', on_delete=models.CASCADE, related_name='stability')
    days = models.ForeignKey('Detail', on_delete=models.CASCADE, related_name='days')
    target_systems = models.ManyToManyField('Asset_System', related_name='target_systems')

class Asset_System(BaseAsset):
    details = models.ManyToManyField('Detail', related_name='details')
    star_clusters = models.ManyToManyField(Asset_Star_Cluster, related_name='star_clusters')
    routes = models.ManyToManyField(Asset_Route, related_name='routes')

class Asset_Sector(BaseAsset):
    systems = models.ManyToManyField(Asset_System, related_name='systems')
