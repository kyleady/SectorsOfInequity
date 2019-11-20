from django.db import models
from django.forms.models import model_to_dict
import json

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

class Region(BaseConfig):
    pass

class WeightedRegion(models.Model):
    weight = models.SmallIntegerField()
    region = models.ForeignKey(Region, on_delete=models.CASCADE)

    def __str__(self):
        return "({weight}) {region_name}".format(weight=self.weight, region_name=self.region.name)

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

class Grid(BaseConfig):
    height = models.PositiveSmallIntegerField(default=20, blank=True)
    width = models.PositiveSmallIntegerField(default=20, blank=True)
    connectionRange =models.PositiveSmallIntegerField(default=5, blank=True)
    populationRate = models.FloatField(default=0.5, blank=True)
    connectionRate = models.FloatField(default=0.5, blank=True)
    rangeRateMultiplier = models.FloatField(default=0.5, blank=True)
    weightedRegions = models.ManyToManyField(WeightedRegion, blank=True)

class Sector(BaseConfig):
    name = models.CharField(default="-", max_length=39)

class SectorSystem(models.Model):
    sector = models.ForeignKey(Sector, on_delete=models.CASCADE)
    x = models.PositiveSmallIntegerField()
    y = models.PositiveSmallIntegerField()
    region = models.ForeignKey(Region, on_delete=models.CASCADE)

class SectorRoute(models.Model):
    start = models.ForeignKey(SectorSystem, on_delete=models.CASCADE, related_name='start')
    end = models.ForeignKey(SectorSystem, on_delete=models.CASCADE, related_name='end')
