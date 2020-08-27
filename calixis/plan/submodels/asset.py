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

    name = models.CharField(max_length=200)

class Asset_Group(BaseAsset):
    type = models.ForeignKey('Config_Name', on_delete=models.CASCADE, related_name='asset_group_type')
    assets = models.ManyToManyField('Asset', related_name='asset_group_assets')

class Asset(BaseAsset):
    def get_details_by_table(self):
        details = self.details.all()
        by_table = {}
        for detail in details:
            inspiration_table = detail.get_inspiration_table_name()
            by_table[inspiration_table] = by_table.get(inspiration_table, [])
            by_table[inspiration_table].append(detail)

        return by_table

    type =  models.ForeignKey('Config_Name', on_delete=models.CASCADE, related_name='asset_type')
    details = models.ManyToManyField('Detail', related_name='asset_details')
    asset_groups = models.ManyToManyField('Asset_Group', related_name='asset_asset_groups')
    grids = models.ManyToManyField('Asset_Grid', related_name='asset_grids')

class Asset_Grid(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(max_length=200)

class Asset_Node(models.Model):
    grid = models.ForeignKey(Asset_Grid, on_delete=models.CASCADE, related_name='grid_node_grid')
    asset = models.ForeignKey('Asset', on_delete=models.CASCADE, related_name='grid_node_asset')
    x = models.PositiveIntegerField()
    y = models.PositiveIntegerField()

class Asset_Connection(models.Model):
    grid = models.ForeignKey(Asset_Grid, on_delete=models.CASCADE, related_name='grid_connection_grid')
    asset = models.ForeignKey('Asset', on_delete=models.CASCADE, related_name='grid_connection_asset')
    start = models.ForeignKey(Asset_Node, on_delete=models.CASCADE, related_name='grid_connection_start')
    end = models.ForeignKey(Asset_Node, on_delete=models.CASCADE, related_name='grid_connection_end')
