from django.urls import path
from .views.default import DefaultViews
from .views.sector_asset import AssetSectorViews
from .views.sector_grid import GridSectorViews
from .views.app_views import AppViews
from .views.subapp_views import SubappViews
from .models import *

app = 'plan'
subapp = 'asset'
asset_models = [
    { 'full_name': 'Sector', 'app': app, 'subapp': subapp, 'name': 'sector', 'Model': Asset_Sector,
        'custom': { 'Config': Grid_Sector, 'Job': Job }, 'Views': AssetSectorViews
    },
    { 'full_name': 'System', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Asset_System },
    { 'full_name': 'Star Cluster', 'app': app, 'subapp': subapp, 'name': 'star-cluster', 'Model': Asset_Star_Cluster },
    { 'full_name': 'Route', 'app': app, 'subapp': subapp, 'name': 'route', 'Model': Asset_Route },
    { 'full_name': 'Zone', 'app': app, 'subapp': subapp, 'name': 'zone', 'Model': Asset_Zone },
]

subapp = 'config'
config_models = [
    { 'full_name': 'Zone Config', 'app': app, 'subapp': subapp, 'name': 'zone', 'Model': Config_Zone },
    { 'full_name': 'Route Config',  'app': app, 'subapp': subapp, 'name': 'route',  'Model': Config_Route },
    { 'full_name': 'Star Cluster Config', 'app': app, 'subapp': subapp, 'name': 'star-cluster', 'Model': Config_Star_Cluster },
    { 'full_name': 'System Config', 'app': app, 'subapp': subapp, 'name': 'system', 'Model': Config_System },
    { 'full_name': 'Perterbation', 'app': app, 'subapp': subapp, 'name': 'perterbation', 'Model': Perterbation },
    { 'full_name': 'Grid Config',   'app': app, 'subapp': subapp, 'name': 'grid',   'Model': Config_Grid },
    { 'full_name': 'Sector Config', 'app': app, 'subapp': subapp, 'name': 'sector', 'Model': Grid_Sector,
        'custom': { 'SubModels': [Grid_System, Grid_Route], 'Grid': Config_Grid, 'Job': Job  }, 'Views': GridSectorViews
    },
]

subapp = 'detail'
detail_models = [
    { 'full_name': 'Detail', 'app': app, 'subapp': subapp, 'name': 'detail', 'Model': Detail },
]

subapp = 'inspiration'
inspiration_models = [
    { 'full_name': 'Inspiration', 'app': app, 'subapp': subapp, 'name': 'inspiration', 'Model': Inspiration },
]


subapp = 'roll'
roll_models = [
    { 'full_name': 'Roll', 'app': app, 'subapp': subapp, 'name': 'roll', 'Model': Roll },
]

subapp = 'weighted'
weighted_models = [
    { 'full_name': 'Weighted Perterbation', 'app': app, 'subapp': subapp, 'name': 'perterbation', 'Model': Weighted_Perterbation },
    { 'full_name': 'Weighted Inspiration', 'app': app, 'subapp': subapp, 'name': 'inspiration', 'Model': Weighted_Inspiration },
]

subapp = 'job'
job_models = [
    { 'full_name': 'Job', 'app': app, 'subapp': subapp, 'name': 'job', 'Model': Job },
]



all_models = (
    asset_models +
    config_models +
    detail_models +
    inspiration_models +
    roll_models +
    weighted_models +
    job_models)

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
                'job'
            ])

urlpatterns.append(
    path(
        '',
        app_views.app,
        name=app_views.app_url
    )
)

all_subapp_views = [
                SubappViews(title='Assets', name=asset_models[0]['subapp'], app=asset_models[0]['app'], models=[
                    model['name'] for model in asset_models
                ]),
                SubappViews(title='Config',name=config_models[0]['subapp'], app=config_models[0]['app'], models=[
                    model['name'] for model in config_models
                ]),
                SubappViews(title='Details', name=detail_models[0]['subapp'], app=detail_models[0]['app'], models=[
                    model['name'] for model in detail_models
                ]),
                SubappViews(title='Inspiration', name=inspiration_models[0]['subapp'], app=inspiration_models[0]['app'], models=[
                    model['name'] for model in inspiration_models
                ]),
                SubappViews(title='Rolls', name=roll_models[0]['subapp'], app=roll_models[0]['app'], models=[
                    model['name'] for model in roll_models
                ]),
                SubappViews(title='Weighted', name=weighted_models[0]['subapp'], app=weighted_models[0]['app'], models=[
                    model['name'] for model in weighted_models
                ]),
                SubappViews(title='Job', name=job_models[0]['subapp'], app=job_models[0]['app'], models=[
                    model['name'] for model in job_models
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
