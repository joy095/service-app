FROM python:3.10-slim

WORKDIR /app/imageCheck

# Install dependencies
COPY ./imageCheck/requirements.txt .
RUN pip install --upgrade pip && pip install -r requirements.txt

# Copy the FastAPI app
COPY ./imageCheck/app .

# Expose the FastAPI port
EXPOSE 8083

# Run the FastAPI app using uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8083"]