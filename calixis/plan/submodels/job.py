from django.db import models
from django.forms.models import model_to_dict
import json

from .config import Config_Grid
from .grid import Grid_Sector

class Job(models.Model):
    class JobType:
        GRID = 'GD'
        SECTOR = 'SR'

        choices = (
            (GRID, 'Grid'),
            (SECTOR, 'Sector'),
        )

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        status = ""
        if self.error:
            status = "ERROR"
        elif 0 <= self.percent_complete and self.asset_id == None:
            status = "{percent_complete}%".format(
                percent_complete=self.percent_complete
            )
        else:
            status = "COMPLETE"

        return "[{status}] {type} {timestamp}".format(
            status=status,
            type=self.jobType,
            timestamp=self.created_at,
        )


    percent_complete = models.PositiveSmallIntegerField(blank=True, null=True, default=0)
    error = models.TextField(blank=True, null=True, default=None)
    created_at = models.DateTimeField(auto_now_add=True)
    jobType = models.CharField(max_length=2, choices=JobType.choices)
    config_id = models.BigIntegerField(null=True)
    asset_id = models.BigIntegerField(null=True, blank=True)
