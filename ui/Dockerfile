FROM debian:latest
RUN apt-get update \
    && apt-get install -y \
        curl \
        git \
        wget \
        zip \
        unzip \
        libgconf-2-4 \
        gdb \
        libstdc++6 \
        libglu1-mesa \
        fonts-droid-fallback \
        lib32stdc++6 \
        python3 \
    && apt-get clean

RUN git clone https://github.com/flutter/flutter.git /usr/local/flutter
ENV PATH="/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin:${PATH}"
RUN flutter doctor -v

# Enable flutter web
RUN flutter channel master
RUN flutter upgrade
RUN flutter config --enable-web

WORKDIR /app
COPY pubspec.yaml .
COPY pubspec.lock .
RUN flutter pub get

COPY . .
RUN flutter create . && flutter build web

FROM nginx:1.17.9-alpine
COPY --from=deps /app/build/web /usr/share/nginx/html
