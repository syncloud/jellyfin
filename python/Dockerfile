FROM python:3.12-slim-bookworm
COPY requirements.txt /
RUN pip install -r /requirements.txt
RUN rm -rf /var/lib/apt/lists/*
