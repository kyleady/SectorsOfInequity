from django.urls import path
from .views.default import DefaultViews
from .views.sector import SectorViews
from .models import *

app = 'plan'
subapp = 'config'
config_models = [
    { 'full_name': 'System Config', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Config_System },
    { 'full_name': 'Region Config', 'app': app, 'subapp': subapp, 'name': 'region', 'Model': Config_Region },
    { 'full_name': 'Grid Config',   'app': app, 'subapp': subapp, 'name': 'grid',   'Model': Config_Grid },
    { 'full_name': 'Sector Config', 'app': app, 'subapp': subapp, 'name': 'sector', 'Model': Config_Sector,
        'custom': { 'SubModels': [Config_Sector_System, Config_Sector_Route], 'Grid': Config_Grid }, 'Views': SectorViews
    },
]

subapp = 'inspiration'
inspiration_models = [
    { 'full_name': 'System Inspiration', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Inspiration_System },
]

subapp = 'perterbation'
perterbation_models = [
    { 'full_name': 'System Perterbataion', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Perterbation_System },
]

subapp = 'weighted'
weighted_models = [
    { 'full_name': 'Weighted Region Config', 'app': app, 'subapp': subapp, 'name': 'config-region', 'Model': Weighted_Config_Region },
    { 'full_name': 'Weighted System Inspiration', 'app': app, 'subapp': subapp, 'name': 'inspiration-system', 'Model': Weighted_Inspiration_System },
]

all_models = config_models + inspiration_models + perterbation_models + weighted_models


{ 'full_name': 'Weighted Region', 'app': app, 'name': 'config-region', 'Model': Weighted_Config_Region },
urlpatterns = []
for model in all_models:
    if 'Views' in model:
        CustomViews = model['Views']
        del model['Views']
        views = CustomViews(**model)
    else:
        views = DefaultViews(**model)

    urlpatterns.append(
        path(
            '{subapp}/{name}/'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.index,
            name=views.index_url
        )
    )
    urlpatterns.append(
        path(
            '{subapp}/{name}/<int:model_id>/'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.detail,
            name=views.detail_url
        )
    )
    urlpatterns.append(
        path(
            '{subapp}/{name}/new/'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.new,
            name=views.new_url
        )
    )
    urlpatterns.append(
        path(
            '{subapp}/{name}/delete/<int:model_id>'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.delete,
            name=views.delete_url
        )
    )
