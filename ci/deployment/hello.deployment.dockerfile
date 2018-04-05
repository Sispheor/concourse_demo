# BUILD:
# docker build --force-rm=true -f hello.deployment.dockerfile -t hello-deployment .

# RUN:
# docker run -it hello-deployment

FROM python:3

RUN pip install docker

ADD deployment.py /deployment.py
RUN chmod +x /deployment.py
CMD ["python /deployment.py"]