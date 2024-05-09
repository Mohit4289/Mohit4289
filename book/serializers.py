from dataclasses import Field
from rest_framework import serializers
from .models import *

class Bookserializer(serializers.ModelSerializer):
    class Meta:
        model = Book
        fields = "__all__"

    