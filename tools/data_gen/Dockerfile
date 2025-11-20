FROM python:3.10-slim

WORKDIR /app

RUN pip install tigramite pandas numpy

COPY gen_data.py .

CMD ["python", "gen_data.py"]