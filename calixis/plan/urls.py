from django.urls import path
from .views.default import DefaultViews
from .views.sector import SectorViews
from .views.app_views import AppViews
from .views.subapp_views import SubappViews
from .models import *

app = 'plan'
subapp = 'asset'
asset_models = [
    { 'full_name': 'Sector', 'app': app, 'subapp': subapp, 'name': 'sector', 'Model': Asset_Sector },
    { 'full_name': 'System', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Asset_System },
]

subapp = 'config'
config_models = [
    { 'full_name': 'System Config', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Config_System },
    { 'full_name': 'Region Config', 'app': app, 'subapp': subapp, 'name': 'region', 'Model': Config_Region },
    { 'full_name': 'Grid Config',   'app': app, 'subapp': subapp, 'name': 'grid',   'Model': Config_Grid },
    { 'full_name': 'Sector Config', 'app': app, 'subapp': subapp, 'name': 'sector', 'Model': Config_Sector,
        'custom': { 'SubModels': [Config_Sector_System, Config_Sector_Route], 'Grid': Config_Grid }, 'Views': SectorViews
    },
]

subapp = 'detail'
detail_models = [
    { 'full_name': 'System Feature Detail', 'app': app, 'subapp': subapp, 'name': 'system-feature', 'Model': Detail_System_Feature },
]

subapp = 'inspiration'
inspiration_models = [
    { 'full_name': 'System Feature Inspiration', 'app': app, 'subapp': subapp, 'name': 'system-feature', 'Model': Inspiration_System_Feature },
]

subapp = 'perterbation'
perterbation_models = [
    { 'full_name': 'System Perterbataion', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Perterbation_System },
]

subapp = 'roll'
roll_models = [
    { 'full_name': 'System Feature Count', 'app': app, 'subapp': subapp, 'name': 'system-features', 'Model': Roll_System_Features },
    { 'full_name': 'System Star Count', 'app': app, 'subapp': subapp, 'name': 'system-stars', 'Model': Roll_System_Stars },
]

subapp = 'weighted'
weighted_models = [
    { 'full_name': 'Weighted Region Config', 'app': app, 'subapp': subapp, 'name': 'config-region', 'Model': Weighted_Config_Region },
    { 'full_name': 'Weighted System Inspiration', 'app': app, 'subapp': subapp, 'name': 'inspiration-system', 'Model': Weighted_Inspiration_System },
]



all_models = (
    asset_models +
    config_models +
    detail_models +
    inspiration_models +
    perterbation_models +
    roll_models +
    weighted_models)

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

app_views = AppViews(title='Sector Planning', name='plan', subapps=[
                'asset',
                'config',
                'detail',
                'weighted',
                'inspiration',
                'roll',
                'perterbation',
            ])

urlpatterns.append(
    path(
        '',
        app_views.app,
        name=app_views.app_url
    )
)

all_subapp_views = [
                SubappViews(title='Assets', name='asset', app='plan', models=[
                    'sector',
                    'system',
                ]),
                SubappViews(title='Config', name='config', app='plan', models=[
                    'system',
                    'region',
                    'grid',
                    'sector',
                ]),
                SubappViews(title='Details', name='detail', app='plan', models=[
                    'system-feature'
                ]),
                SubappViews(title='Inspiration', name='inspiration', app='plan', models=[
                    'system-feature',
                ]),
                SubappViews(title='Perterbations', name='perterbation', app='plan', models=[
                    'system',
                ]),
                SubappViews(title='Rolls', name='roll', app='plan', models=[
                    'system-features',
                    'system-stars',
                ]),
                SubappViews(title='Weighted', name='weighted', app='plan', models=[
                    'config-region',
                    'inspiration-system',
                ]),
]

for subapp_views in all_subapp_views:
    urlpatterns.append(
        path(
            '{subapp}'.format(subapp=subapp_views.name),
            subapp_views.subapp,
            name=subapp_views.subapp_url
        )
    )
