from django.shortcuts import get_object_or_404, render, redirect
from django import forms
import requests
import json

from .default import DefaultViews

class AssetSectorViews(DefaultViews):
    def new(self, request):
        allConfig = self.custom['Config'].objects
        class GenerateSectorAssetForm(forms.Form):
            config = forms.ModelChoiceField(queryset=allConfig)
        if request.POST:
            form = GenerateSectorAssetForm(request.POST)
            form.is_valid()
            config = form.cleaned_data['config']
            Job = self.custom['Job']
            jobType = Job.JobType.SECTOR
            job = Job(jobType=jobType, config_id=config.id)
            job.save()
            screaming_vortex_url = 'http://{host}/{path}'.format(
                host=self.screaming_vortex_host,
                path='sector'
            )
            screaming_vortex_json = {
                'config_id': config.id,
                'job_id': job.id,
            }
            screaming_vortex_response = requests.post(
                screaming_vortex_url,
                data=json.dumps(screaming_vortex_json),
            )
            screaming_vortex_response.raise_for_status()
            return redirect(self.index_url)
        else:
            template = 'detail.html'
            context = {
                'full_name': self.full_name,
                'new_url': self.new_url,
                'detail_url': self.detail_url,
                'form': GenerateSectorAssetForm()
            }

            return render(request, template, context)
