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

class Asset_Sector(BaseAsset):
    pass

class Asset_System(BaseAsset):
    sector = models.ForeignKey(Asset_Sector, on_delete=models.CASCADE)
