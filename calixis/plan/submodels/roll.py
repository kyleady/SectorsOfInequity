from django.db import models
from django.forms.models import model_to_dict
from django.core.validators import int_list_validator
import json

class Roll(models.Model):
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
            text += "kl{drop_highest}".format(
                        drop_highest=-1*self.keep_highest
                    )

        if self.multiplier != 1 and text != "":
            text = "{multiplier}x({text})".format(
                        multiplier=self.multiplier,
                        text=text
                    )

        if self.base > 0 and text != "":
            text += "+"

        if self.base != 0:
            text += "{base}".format(
                        base=self.base
                    )

        if text == "":
            text = "0"

        if self.minimum:
            text = "min({text}, {min})".format(
                text=text,
                min=self.minimum
            )

        if self.maximum:
            text = "max({text}, {max})".format(
                text=text,
                min=self.maximum
            )

        return text

    @classmethod
    def parse(cls, roll_str):
        roll_match = re.match("(?:(\d*)(?:d|D)(\d+))?(?:(kh|dh)(\d+))?(+|-|)(\d*)", roll_str)

        if not roll_match:
            return None

        dice_count = 0
        dice_size = 0
        keep_highest = 0
        base = 0
        multiplier = 1

        if roll_match[1]:
            dice_size = int(roll_match[1])
            if roll_match[0]:
                dice_count = int(roll_match[0])
            else:
                dice_count = 1

        if roll_match[2]:
            if roll_match[2] == "kh":
                keep_highest = int(roll_match[3])
            elif roll_match[2] == "kl":
                keep_highest = -1 * int(roll_match[3])

        if roll_match[5]:
            if roll_match[4] != "-":
                base = int(roll_match[5])
            else:
                base = -1 * int(roll_match[5])

        return cls(
            dice_count=dice_count,
            dice_size=dice_size,
            keep_highest=keep_highest,
            base=base,
            multiplier=multiplier,
        )

    required_flags = models.CharField(blank=True, null=True, max_length=200)
    rejected_flags = models.CharField(blank=True, null=True, max_length=200)
    dice_count = models.PositiveSmallIntegerField(blank=True, default=0)
    dice_size = models.PositiveSmallIntegerField(blank=True, default=6)
    base = models.IntegerField(blank=True, default=0)
    multiplier = models.IntegerField(blank=True, default=1)
    keep_highest = models.IntegerField(blank=True, default=0)
    minimum = models.PositiveSmallIntegerField(blank=True, null=True)
    maximum = models.PositiveSmallIntegerField(blank=True, null=True)
    rolls = models.ManyToManyField('Roll', related_name='roll_rolls')
