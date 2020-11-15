FROM python:3.6.1-alpine
RUN apk add --no-cache python3-dev \
    && pip3 install --upgrade pip
WORKDIR /app
ADD . /app
RUN pip3 --no-cache-dir install -r requeriments.txt
EXPOSE 82
CMD ["python3","receptor.py"]
