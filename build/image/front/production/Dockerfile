FROM node:10
WORKDIR /front
RUN git clone https://github.com/frontend-park-mail-ru/2019_1_Escapade.git . && \ 
    git checkout develop
RUN npm install
COPY configuration/external/front/production/ .
RUN npm run build
ENTRYPOINT ["npm"]