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

    name = models.CharField(default="-", max_length=25)

# SystemConfig Model
class Config_System(BaseConfig):
    pass

# RegionConfig Model
class Config_Region(BaseConfig):
    system = models.ForeignKey(Config_System, on_delete=models.CASCADE)

# GridConfig Model
class Config_Grid(BaseConfig):
    height = models.PositiveSmallIntegerField(default=20, blank=True)
    width = models.PositiveSmallIntegerField(default=20, blank=True)
    connectionRange =models.PositiveSmallIntegerField(default=5, blank=True)
    populationRate = models.FloatField(default=0.5, blank=True)
    connectionRate = models.FloatField(default=0.5, blank=True)
    rangeRateMultiplier = models.FloatField(default=0.5, blank=True)
    smoothingFactor = models.FloatField(default=0.5, blank=True)

# SectorConfig Model
class Config_Sector(BaseConfig):
    name = models.CharField(default="-", max_length=39)

class Config_Sector_System(models.Model):
    sector = models.ForeignKey(Config_Sector, on_delete=models.CASCADE)
    x = models.PositiveSmallIntegerField()
    y = models.PositiveSmallIntegerField()
    region = models.ForeignKey(Config_Region, on_delete=models.CASCADE)

class Config_Sector_Route(models.Model):
    start = models.ForeignKey(Config_Sector_System, on_delete=models.CASCADE, related_name='start')
    end = models.ForeignKey(Config_Sector_System, on_delete=models.CASCADE, related_name='end')
