from django.db import models
from django.forms.models import model_to_dict
from django.core.validators import int_list_validator
import json
import re

from .asset import Asset_System
from .inspiration import Inspiration_System_Feature

# Abstract Models
class BaseDetail(models.Model):
    class Meta:
        abstract = True

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        roll_list = self.rolls.split(',')
        text = self.inspiration.description
        for roll in roll_list:
            text = re.sub(r'\[\[[^\]]\]\]', roll, text)
        return text

    rolls = models.CharField(validators=[int_list_validator], max_length=100)
    asset = None
    inspiration = None


class Detail_System_Feature(BaseDetail):
    asset = models.ForeignKey(Asset_System, on_delete=models.CASCADE)
    inspiration = models.ForeignKey(Inspiration_System_Feature, on_delete=models.CASCADE)
