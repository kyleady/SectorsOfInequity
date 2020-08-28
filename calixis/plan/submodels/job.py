from django.db import models
from django.forms.models import model_to_dict
import json

class Job(models.Model):
    class JobType:
        ASSET = 'AS'
        OTHER = 'OT'

        choices = (
            (ASSET, 'Asset'),
            (OTHER, 'Other')
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

    percent_complete = models.PositiveSmallIntegerField(blank=True, default=0)
    error = models.TextField(blank=True, null=True)
    started_at = models.DateTimeField(auto_now_add=True)
    finised_at = models.DateTimeField(blank=True, null=True)
    type = models.ForeignKey('Config_Name', null=True, on_delete=models.SET_NULL, related_name='job_type')
    perterbation = models.ForeignKey('Perterbation', null=True, on_delete=models.SET_NULL, related_name='job_perterbation')
    asset = models.ForeignKey('Asset', null=True, blank=True, on_delete=models.SET_NULL, related_name='job_asset')
