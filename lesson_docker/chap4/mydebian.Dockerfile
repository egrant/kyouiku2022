FROM debian:bullseye

RUN apt update && apt install -y \
  bpytop \
  && rm -rf /var/lib/apt/lists/*

RUN apt update && apt install -y \
  emacs \
  git \
  less \
  ssh \
  vim \
  wget \
  && rm -rf /var/lib/apt/lists/*

RUN apt update && apt install -y \
  zsh \
  && rm -rf /var/lib/apt/lists/*

CMD ["zsh"]
