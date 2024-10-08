ARG VARIANT=bookworm-slim
FROM debian:${VARIANT} as nvim
ARG NVIM_VERSION=0.10.1
WORKDIR /tmp
RUN apt update && apt install -y git cmake gettext
RUN git clone https://github.com/neovim/neovim.git --branch "v${NVIM_VERSION}" --depth=1 && \
    cd neovim && \
    make CMAKE_BUILD_TYPE=RelWithDebInfo && \
    make install

FROM debian:${VARIANT}
COPY --from=nvim /usr/local/bin/nvim /usr/local/bin/nvim
COPY --from=nvim /usr/local/lib/nvim /usr/local/lib/nvim
COPY --from=nvim /usr/local/share/nvim /usr/local/share/nvim
ADD ./nvim-configs /usr/local/share/nvim/configs


# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
