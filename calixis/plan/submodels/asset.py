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

    name = models.CharField(default="-", max_length=25)
    parent = None


class Asset_Star_Cluster(BaseAsset):
    stars = models.ManyToManyField('Detail', related_name='stars')

class Asset_Route(BaseAsset):
    type = models.ForeignKey('Detail', on_delete=models.CASCADE, related_name='type')
    days = models.ForeignKey('Detail', on_delete=models.CASCADE, related_name='days')
    target_systems = models.ManyToManyField('Asset_System', related_name='target_systems')

class Asset_System(BaseAsset):
    details = models.ManyToManyField('Detail', related_name='details')
    star_clusters = models.ManyToManyField(Asset_Star_Cluster, related_name='star_clusters')
    routes = models.ManyToManyField(Asset_Route, related_name='routes')

class Asset_Sector(BaseAsset):
    systems = models.ManyToManyField(Asset_System, related_name='systems')
