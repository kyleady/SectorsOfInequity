from django.db import models

class Tag(models.Model):
    def __str__(self):
        return self.name

    name = models.CharField(max_length=100)
