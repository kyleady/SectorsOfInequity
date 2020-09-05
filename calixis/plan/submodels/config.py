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

    name = models.CharField(max_length=200)

class Tag(BaseConfig):
    pass

class Config_Name(BaseConfig):
    tags =  models.ManyToManyField(Tag, related_name='config_type_tags')

class Config_Asset(BaseConfig):
    type = models.ForeignKey('Config_Name', on_delete=models.CASCADE, related_name='config_asset_type')
    order = models.ManyToManyField('Roll', related_name='config_asset_order')
    inspiration_tables = models.ManyToManyField('Weighted_Table', related_name='config_asset_inspiration_tables')
    child_configs = models.ManyToManyField('Config_Group', related_name='config_asset_child_configs')
    grids = models.ManyToManyField('Config_Grid', related_name='config_asset_grids')
    tags = models.ManyToManyField(Tag, related_name='config_asset_tags')

class Config_Group(BaseConfig):
    types = models.ManyToManyField('Weighted_Type', related_name='config_sub_asset_types')
    count = models.ManyToManyField('Roll', related_name='config_sub_asset_count')
    extras = models.ManyToManyField('Inspiration_Extra', related_name='config_sub_asset_extras')
    perterbations = models.ManyToManyField('Perterbation', related_name='config_sub_asset_perterbations')
    tags = models.ManyToManyField(Tag, related_name='config_sub_asset_tags')

class Config_Region(BaseConfig):
    types = models.ManyToManyField('Weighted_Type', related_name='config_region_types')
    perterbations = models.ManyToManyField('Perterbation', related_name='config_region_perterbations')
    tags = models.ManyToManyField(Tag, related_name='config_region_tags')

class Config_Grid(BaseConfig):
    regions = models.ManyToManyField('Weighted_Region', related_name='config_grid_regions')
    connection_types = models.ManyToManyField('Weighted_Type', related_name='config_grid_connection_types')
    count = models.ManyToManyField('Roll', related_name='config_grid_count')
    height = models.ManyToManyField('Roll', related_name='config_grid_height')
    width = models.ManyToManyField('Roll', related_name='config_grid_width')
    connection_range = models.ManyToManyField('Roll', related_name='config_grid_connection_range')
    population_percent = models.ManyToManyField('Roll', related_name='config_grid_population_percent')
    population_denominator = models.IntegerField(default=100, blank=True)
    connection_percent = models.ManyToManyField('Roll', related_name='config_grid_connection_percent')
    connection_denominator = models.IntegerField(default=100, blank=True)
    range_multiplier_percent =  models.ManyToManyField('Roll', related_name='config_grid_range_multiplier_percent')
    range_multiplier_denominator = models.IntegerField(default=100, blank=True)
    smoothing_percent = models.ManyToManyField('Roll', related_name='config_grid_smoothing_percent')
    smoothing_denominator = models.IntegerField(default=100, blank=True)
    tags = models.ManyToManyField(Tag, related_name='config_grid_tags')
