from django.db import models
import uuid
from django.contrib.auth.models import User

# Create your models here.
class Book(models.Model):
    user = models.ForeignKey(User, on_delete=models.SET_NULL, null=True, blank=True)
    title = models.CharField(max_length=20)
    author = models.CharField(max_length=20)
    description = models.TextField(max_length=50)
    published = models.DateField()

    def __str__(self):
        return self.title
    