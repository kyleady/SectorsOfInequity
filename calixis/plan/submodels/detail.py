from django.db import models
from django.forms.models import model_to_dict
from django.core.validators import int_list_validator
import json
import re

class Detail(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.inspirations.all()[0].name

    def get_rolls(self):
        return self.rolls.split(',')

    def get_nested_name(self):
        nested_inspirations = self.nested_inspirations.all()
        if len(self.nested_inspirations.all()) > 0:
            return nested_inspirations[0].name
        else:
            return "-"

    def get_description(self):
        text = '\n'.join(inspiration.description for inspiration in self.inspirations.all())
        for roll in self.get_rolls():
            text = re.sub(r'\[\[\d+\]\]', roll, text, count=1)
        return text

    rolls = models.CharField(validators=[int_list_validator], max_length=100)
    inspirations = models.ManyToManyField('Inspiration', related_name='inspirations')
    nested_inspirations = models.ManyToManyField('Inspiration_Nested')
    parent_detail = models.ForeignKey('self', null=True, blank=True, on_delete=models.CASCADE)
