# Generated by Django 5.0.6 on 2024-05-09 07:50

import uuid
from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('book', '0001_initial'),
    ]

    operations = [
        migrations.AlterField(
            model_name='book',
            name='id',
            field=models.UUIDField(default=uuid.uuid4, primary_key=True, serialize=False),
        ),
    ]