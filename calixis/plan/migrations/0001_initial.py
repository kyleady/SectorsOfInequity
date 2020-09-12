# Generated by Django 2.2.9 on 2020-09-10 18:01

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Asset',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Asset_Grid',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
            ],
        ),
        migrations.CreateModel(
            name='Config_Asset',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Config_Name',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Config_Region',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Inspiration',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=75)),
                ('description', models.TextField()),
                ('extra_rolls', models.SmallIntegerField(blank=True, default=0)),
            ],
        ),
        migrations.CreateModel(
            name='Inspiration_Table',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=100)),
            ],
        ),
        migrations.CreateModel(
            name='Roll',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('required_flags', models.CharField(blank=True, max_length=200, null=True)),
                ('rejected_flags', models.CharField(blank=True, max_length=200, null=True)),
                ('dice_count', models.PositiveSmallIntegerField(blank=True, default=0)),
                ('dice_size', models.PositiveSmallIntegerField(blank=True, default=6)),
                ('base', models.IntegerField(blank=True, default=0)),
                ('multiplier', models.IntegerField(blank=True, default=1)),
                ('keep_highest', models.IntegerField(blank=True, default=0)),
                ('minimum', models.PositiveSmallIntegerField(blank=True, null=True)),
                ('maximum', models.PositiveSmallIntegerField(blank=True, null=True)),
                ('rolls', models.ManyToManyField(related_name='roll_rolls', to='plan.Roll')),
            ],
        ),
        migrations.CreateModel(
            name='Tag',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Weighted_Type',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('order', models.ManyToManyField(related_name='_weighted_type_order_+', to='plan.Roll')),
                ('value', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_Name')),
                ('weights', models.ManyToManyField(related_name='_weighted_type_weights_+', to='plan.Roll')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Weighted_Table',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('order', models.ManyToManyField(related_name='_weighted_table_order_+', to='plan.Roll')),
                ('value', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Inspiration_Table')),
                ('weights', models.ManyToManyField(related_name='_weighted_table_weights_+', to='plan.Roll')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Weighted_Region',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('order', models.ManyToManyField(related_name='_weighted_region_order_+', to='plan.Roll')),
                ('value', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_Region')),
                ('weights', models.ManyToManyField(related_name='_weighted_region_weights_+', to='plan.Roll')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Weighted_Inspiration',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('order', models.ManyToManyField(related_name='_weighted_inspiration_order_+', to='plan.Roll')),
                ('value', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Inspiration')),
                ('weights', models.ManyToManyField(related_name='_weighted_inspiration_weights_+', to='plan.Roll')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Perterbation',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='', max_length=100)),
                ('flags', models.CharField(blank=True, max_length=100, null=True)),
                ('muted_flags', models.CharField(blank=True, max_length=100, null=True)),
                ('required_flags', models.CharField(blank=True, max_length=200, null=True)),
                ('rejected_flags', models.CharField(blank=True, max_length=200, null=True)),
                ('configs', models.ManyToManyField(related_name='perterbation_configs', to='plan.Config_Asset')),
                ('tags', models.ManyToManyField(related_name='perterbation_tags', to='plan.Tag')),
            ],
        ),
        migrations.CreateModel(
            name='Job',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('percent_complete', models.PositiveSmallIntegerField(blank=True, default=0)),
                ('error', models.TextField(blank=True, null=True)),
                ('started_at', models.DateTimeField(auto_now_add=True)),
                ('finised_at', models.DateTimeField(blank=True, null=True)),
                ('asset', models.ForeignKey(blank=True, null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='job_asset', to='plan.Asset')),
                ('perterbation', models.ForeignKey(null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='job_perterbation', to='plan.Perterbation')),
                ('type', models.ForeignKey(null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='job_type', to='plan.Config_Name')),
            ],
        ),
        migrations.AddField(
            model_name='inspiration_table',
            name='count',
            field=models.ManyToManyField(related_name='inspiration_table_count', to='plan.Roll'),
        ),
        migrations.AddField(
            model_name='inspiration_table',
            name='extra_inspirations',
            field=models.ManyToManyField(related_name='inspiration_table_extra_inspirations', to='plan.Weighted_Inspiration'),
        ),
        migrations.AddField(
            model_name='inspiration_table',
            name='modifiers',
            field=models.ManyToManyField(related_name='inspiration_table_modifiers', to='plan.Roll'),
        ),
        migrations.AddField(
            model_name='inspiration_table',
            name='tags',
            field=models.ManyToManyField(related_name='inspiration_table_tags', to='plan.Tag'),
        ),
        migrations.AddField(
            model_name='inspiration_table',
            name='weighted_inspirations',
            field=models.ManyToManyField(related_name='inspiration_table_weighted_inspirations', to='plan.Weighted_Inspiration'),
        ),
        migrations.CreateModel(
            name='Inspiration_Extra',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=100)),
                ('count', models.ManyToManyField(related_name='extra_tables_count', to='plan.Roll')),
                ('inspiration_tables', models.ManyToManyField(related_name='inspiration_tables', to='plan.Weighted_Table')),
                ('tags', models.ManyToManyField(related_name='inspiration_extra_tags', to='plan.Tag')),
                ('type', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='inspiration_extra_type', to='plan.Config_Name')),
            ],
        ),
        migrations.AddField(
            model_name='inspiration',
            name='inspiration_tables',
            field=models.ManyToManyField(related_name='sub_inspiration_tables', to='plan.Weighted_Table'),
        ),
        migrations.AddField(
            model_name='inspiration',
            name='perterbations',
            field=models.ManyToManyField(to='plan.Perterbation'),
        ),
        migrations.AddField(
            model_name='inspiration',
            name='roll_groups',
            field=models.ManyToManyField(related_name='roll_groups', to='plan.Inspiration_Table'),
        ),
        migrations.AddField(
            model_name='inspiration',
            name='tags',
            field=models.ManyToManyField(related_name='inspiration_tags', to='plan.Tag'),
        ),
        migrations.CreateModel(
            name='Detail',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('rolls', models.CharField(max_length=100)),
                ('inspiration_tables', models.ManyToManyField(to='plan.Inspiration_Table')),
                ('inspirations', models.ManyToManyField(related_name='inspirations', to='plan.Inspiration')),
                ('parent_detail', models.ForeignKey(blank=True, null=True, on_delete=django.db.models.deletion.CASCADE, to='plan.Detail')),
            ],
        ),
        migrations.AddField(
            model_name='config_region',
            name='perterbations',
            field=models.ManyToManyField(related_name='config_region_perterbations', to='plan.Perterbation'),
        ),
        migrations.AddField(
            model_name='config_region',
            name='tags',
            field=models.ManyToManyField(related_name='config_region_tags', to='plan.Tag'),
        ),
        migrations.AddField(
            model_name='config_region',
            name='types',
            field=models.ManyToManyField(related_name='config_region_types', to='plan.Weighted_Type'),
        ),
        migrations.AddField(
            model_name='config_name',
            name='tags',
            field=models.ManyToManyField(related_name='config_type_tags', to='plan.Tag'),
        ),
        migrations.CreateModel(
            name='Config_Group',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
                ('count', models.ManyToManyField(related_name='config_sub_asset_count', to='plan.Roll')),
                ('extras', models.ManyToManyField(related_name='config_sub_asset_extras', to='plan.Inspiration_Extra')),
                ('perterbations', models.ManyToManyField(related_name='config_sub_asset_perterbations', to='plan.Perterbation')),
                ('tags', models.ManyToManyField(related_name='config_sub_asset_tags', to='plan.Tag')),
                ('types', models.ManyToManyField(related_name='config_sub_asset_types', to='plan.Weighted_Type')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Config_Grid',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
                ('population_denominator', models.IntegerField(blank=True, default=100)),
                ('connection_denominator', models.IntegerField(blank=True, default=100)),
                ('range_multiplier_denominator', models.IntegerField(blank=True, default=100)),
                ('smoothing_denominator', models.IntegerField(blank=True, default=100)),
                ('connection_percent', models.ManyToManyField(related_name='config_grid_connection_percent', to='plan.Roll')),
                ('connection_range', models.ManyToManyField(related_name='config_grid_connection_range', to='plan.Roll')),
                ('connection_types', models.ManyToManyField(related_name='config_grid_connection_types', to='plan.Weighted_Type')),
                ('count', models.ManyToManyField(related_name='config_grid_count', to='plan.Roll')),
                ('height', models.ManyToManyField(related_name='config_grid_height', to='plan.Roll')),
                ('population_percent', models.ManyToManyField(related_name='config_grid_population_percent', to='plan.Roll')),
                ('range_multiplier_percent', models.ManyToManyField(related_name='config_grid_range_multiplier_percent', to='plan.Roll')),
                ('regions', models.ManyToManyField(related_name='config_grid_regions', to='plan.Weighted_Region')),
                ('smoothing_percent', models.ManyToManyField(related_name='config_grid_smoothing_percent', to='plan.Roll')),
                ('tags', models.ManyToManyField(related_name='config_grid_tags', to='plan.Tag')),
                ('width', models.ManyToManyField(related_name='config_grid_width', to='plan.Roll')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.AddField(
            model_name='config_asset',
            name='child_configs',
            field=models.ManyToManyField(related_name='config_asset_child_configs', to='plan.Config_Group'),
        ),
        migrations.AddField(
            model_name='config_asset',
            name='grids',
            field=models.ManyToManyField(related_name='config_asset_grids', to='plan.Config_Grid'),
        ),
        migrations.AddField(
            model_name='config_asset',
            name='inspiration_tables',
            field=models.ManyToManyField(related_name='config_asset_inspiration_tables', to='plan.Weighted_Table'),
        ),
        migrations.AddField(
            model_name='config_asset',
            name='order',
            field=models.ManyToManyField(related_name='config_asset_order', to='plan.Roll'),
        ),
        migrations.AddField(
            model_name='config_asset',
            name='tags',
            field=models.ManyToManyField(related_name='config_asset_tags', to='plan.Tag'),
        ),
        migrations.AddField(
            model_name='config_asset',
            name='type',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='config_asset_type', to='plan.Config_Name'),
        ),
        migrations.CreateModel(
            name='Asset_Node',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('x', models.PositiveIntegerField()),
                ('y', models.PositiveIntegerField()),
                ('asset', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='grid_node_asset', to='plan.Asset')),
                ('grid', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='grid_node_grid', to='plan.Asset_Grid')),
            ],
        ),
        migrations.CreateModel(
            name='Asset_Group',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200)),
                ('assets', models.ManyToManyField(related_name='asset_group_assets', to='plan.Asset')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Asset_Connection',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('asset', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='grid_connection_asset', to='plan.Asset')),
                ('end', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='grid_connection_end', to='plan.Asset_Node')),
                ('grid', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='grid_connection_grid', to='plan.Asset_Grid')),
                ('start', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='grid_connection_start', to='plan.Asset_Node')),
            ],
        ),
        migrations.AddField(
            model_name='asset',
            name='asset_groups',
            field=models.ManyToManyField(related_name='asset_asset_groups', to='plan.Asset_Group'),
        ),
        migrations.AddField(
            model_name='asset',
            name='details',
            field=models.ManyToManyField(related_name='asset_details', to='plan.Detail'),
        ),
        migrations.AddField(
            model_name='asset',
            name='grids',
            field=models.ManyToManyField(related_name='asset_grids', to='plan.Asset_Grid'),
        ),
        migrations.AddField(
            model_name='asset',
            name='type',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='asset_type', to='plan.Config_Name'),
        ),
    ]
