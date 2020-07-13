from django.db import models

class BaseWeighted(models.Model):
    class Meta:
        abstract = True

    weights = models.ManyToManyField('Roll')

    def __str__(self):
        return "({weights}) {value_name}".format(weights=self.get_weights_as_str(), value_name=self.value.name)

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

class Weighted_Inspiration(BaseWeighted):
    value = models.ForeignKey('Inspiration', on_delete=models.CASCADE)

class Weighted_Perterbation(BaseWeighted):
    value = models.ForeignKey('Perterbation', on_delete=models.CASCADE)
    parent = models.ForeignKey('Config_Grid', on_delete=models.CASCADE)
