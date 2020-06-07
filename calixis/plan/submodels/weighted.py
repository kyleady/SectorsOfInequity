from django.db import models

class BaseWeighted(models.Model):
    class Meta:
        abstract = True

    weight = models.SmallIntegerField(blank=True, default=1)

    def __str__(self):
        return "({weight}) {value_name}".format(weight=self.weight, value_name=self.value.name)

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
