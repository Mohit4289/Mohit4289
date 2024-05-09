
from django.shortcuts import render,redirect, HttpResponse
from .forms import BookForm
from .serializers import Bookserializer
from rest_framework.response import Response
from rest_framework.decorators import api_view
from .models import *
from .pagination import booklistpagination
from rest_framework.permissions import IsAuthenticated
from rest_framework.authentication import TokenAuthentication



# Create your views here.
def home(request):
    if request.method == 'POST':
        form = BookForm(request.POST)
        if form.is_valid():
            form.save()
            return render(request, 'home.html')
    else:
        form = BookForm()
    return render(request, 'home.html', {'form': form})

@api_view(['GET','POST'])
def booklist(request):
    if request.method == 'GET':
        authentication_classes = [TokenAuthentication]
        permission_classes = [IsAuthenticated]
        bok = Book.objects.all()
        pagination_class = booklistpagination
        serializer = Bookserializer(bok, many=True)
        return Response(serializer.data)
    if request.method == 'POST':
        serializer = Bookserializer(data=request.data)
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data)
        else:
            return Response(serializer.errors)


@api_view(['GET','PUT','DElETE'])
def bookdetail(request,pk):
    if request.method == 'GET':
        bok = Book.objects.get(pk=pk)
        serializer = Bookserializer(bok)
        return Response(serializer.data)
    
    if request.method == 'PUT':
        bok = Book.objects.get(pk=pk)
        serializer = Bookserializer(bok, data=request.data)
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data)
        else:
            return Response(serializer.errors)
        
    if request.method == 'DELETE':
        bok = Book.objects.get(pk=pk)
        bok.delete()
        return HttpResponse(status=204)