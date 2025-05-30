FROM python:3.9-alpine

WORKDIR /app

RUN pip install flask docker

COPY . .

ENV PORT=7777

EXPOSE ${PORT}

CMD ["python", "app.py"]