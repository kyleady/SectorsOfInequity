from django.db import models
from django.forms.models import model_to_dict
from django.core.validators import int_list_validator
import json

from .config import Config_System

# Abstract Models
class BaseRoll(models.Model):
    class Meta:
        abstract = True

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        text = ""
        if self.dice_count > 0 and self.dice_size > 0:
            text += "{dice_count}D{dice_size}".format(
                        dice_count=self.dice_count,
                        dice_size=self.dice_size
                    )
        if self.keep_highest > 0:
            text += "kh{keep_highest}".format(
                        keep_highest=self.keep_highest
                    )
        elif self.keep_highest < 0:
            text += "dh{drop_highest}".format(
                        drop_highest=-1*self.keep_highest
                    )
        if self.base > 0:
            text += "+{base}".format(
                        base=self.base
                    )
        elif self.base < 0:
            text += "{base}".format(
                        base=self.base
                    )
        elif text == "":
            text = "0"
        return text

    def parse(self, roll_str):
        roll_match = re.match("(?:(\d*)d(\d+))?(?:(kh|dh)(\d+))?(+|-|)(\d*)", roll_str)

        if not roll_match:
            return False

        if roll_match[1]:
            dice_size = int(roll_match[1])
            if roll_match[0]:
                dice_count = int(roll_match[0])
            else:
                dice_count = 1
        else:
            dice_count = 0
            dice_size = 0

        if roll_match[2]:
            if roll_match[2] == "kh":
                keep_highest = int(roll_match[3])
            elif roll_match[2] == "dh":
                keep_highest = -1 * int(roll_match[3])
        else:
            keep_highest = 0

        if roll_match[5]:
            if roll_match[4] != "-":
                base = int(roll_match[5])
            else:
                base = -1 * int(roll_match[5])
        else:
            base = 0

        return True

    dice_count = models.PositiveSmallIntegerField()
    dice_size = models.PositiveSmallIntegerField()
    base = models.IntegerField()
    multiplier = models.IntegerField()
    keep_highest = models.IntegerField()


class Roll_System_Features(models.Model):
    system = models.ForeignKey(Config_System, on_delete=models.CASCADE)

class Roll_System_Stars(models.Model):
    system = models.ForeignKey(Config_System, on_delete=models.CASCADE)
