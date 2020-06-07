from django.db import models
from django.forms.models import model_to_dict
from django.core.validators import int_list_validator
import json
import re

from .inspiration import Inspiration

class Detail(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.inspiration.name

    def get_description(self):
        roll_list = self.rolls.split(',')
        text = self.inspiration.description
        for roll in roll_list:
            text = re.sub(r'\[\[[^\]]\]\]', roll, text)
        return text

    rolls = models.CharField(validators=[int_list_validator], max_length=100)
    inspiration = models.ForeignKey(Inspiration, on_delete=models.CASCADE)
    parent_detail = models.ForeignKey('self', null=True, blank=True, on_delete=models.CASCADE)
