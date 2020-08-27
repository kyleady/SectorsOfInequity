from django.db import models
from django.forms.models import model_to_dict

class Weighted_Value(models.Model):
    class Meta:
        abstract = True

    def get_weights_as_str(self):
        weights = []
        has_conditional_weights = False
        for weight in self.weights.all():
            if not weight.required_flags:
                weights.append(str(weight))
            else:
                has_conditional_weights = True

        weights_as_str = "+".join(weights)
        if has_conditional_weights:
            weights_as_str += "*"
        return weights_as_str

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return "({weights}) {value_name}".format(weights=self.get_weights_as_str(), value_name=self.value.name)

    weights = models.ManyToManyField('Roll')
    value = None

class Weighted_Inspiration(Weighted_Value):
    value = models.ForeignKey('Inspiration', on_delete=models.CASCADE)

class Weighted_Region(Weighted_Value):
    value = models.ForeignKey('Config_Region', on_delete=models.CASCADE)
