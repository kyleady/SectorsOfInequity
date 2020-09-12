from django.shortcuts import get_object_or_404, render, redirect
from django import forms
import requests
import json

from .default import DefaultViews

class AssetViews(DefaultViews):
    def new(self, request):
        namedPerterbations = self.custom['Perterbation_Model'].objects.filter(tags__name="_Primary")
        allTypes = self.custom['Config_Name_Model'].objects
        class GenerateAssetForm(forms.Form):
            perterbation = forms.ModelChoiceField(queryset=namedPerterbations)
            type = forms.ModelChoiceField(queryset=allTypes)
        if request.POST:
            form = GenerateAssetForm(request.POST)
            form.is_valid()
            perterbation = form.cleaned_data['perterbation']
            type = form.cleaned_data['type']
            Job = self.custom['Job']
            job = Job(perterbation_id=perterbation.id, type_id=type.id)
            job.save()
            screaming_vortex_url = 'http://{host}/{path}'.format(
                host=self.screaming_vortex_host,
                path='asset'
            )
            screaming_vortex_json = {
                'perterbation_id': perterbation.id,
                'type_id': type.id,
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
                'form': GenerateAssetForm()
            }

            return render(request, template, context)
