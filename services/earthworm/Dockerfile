FROM debian:stable-slim

ARG version

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib/ \
    PYTHONPATH=$PYTHONPATH:/usr/local/lib/python3.5/dist-packages/ \
    YAAFE_PATH=/usr/local/lib/python3.5/dist-packages/yaafelib/yaafe_extensions \
    YAAFE_VERSION="${version:-0.70}"

RUN set -x \
    # setup
    && apt-get update \
    && apt-get install -y \
        cmake \
        curl \
        gcc \
        g++ \
        libargtable2-0 \
        libargtable2-dev \
        libeigen3-dev \
        libfftw3-3 \
        libfftw3-dev \
        libhdf5-dev \
        liblapack-dev \
        libmpg123-0 \
        libmpg123-dev \
        libsndfile1 \
        libsndfile1-dev \
        python3.5 \
        python3.5-dev \
        python3-pip \
    && pip3 install -q \
        numpy \
    # build yaafe
    && cd /tmp \
    && curl -fSLO "https://github.com/Yaafe/Yaafe/archive/${YAAFE_VERSION}.tar.gz" \
    && tar -xzf "${YAAFE_VERSION}.tar.gz" \
    && mkdir "Yaafe-${YAAFE_VERSION}/build" \
    && cd "Yaafe-${YAAFE_VERSION}/build" \
    && cmake \
        -DCMAKE_INSTALL_PREFIX=/usr/local/ \
        -DCMAKE_INSTALL_PYTHON_PACKAGES=/usr/local/lib/python3.5/dist-packages/ \
        -DWITH_FFTW3=ON \
        -DHDF5_ROOT=/usr/lib/x86_64-linux-gnu/hdf5/serial/ \
        -DWITH_HDF5=ON \
        -DWITH_LAPACK=ON -DWITH_MATLAB_MEX=OFF -DWITH_MPG123=ON -DWITH_SNDFILE=ON \
        .. \
    && make \
    && make install \
    # cleanup
    && apt-get remove --purge -y \
        cmake \
        cmake-data \
        cpp \
        cpp-6 \
        curl \
        gcc \
        g++ \
        g++-6 \
        gcc-6 \
        gfortran-6 \
        libargtable2-dev \
        libgcc-6-dev \
        libgfortran-6-dev \
        libmpg123-dev \
        libpython3-dev \
        libpython3.5-dev \
        libsndfile1-dev \
        libstdc++-6-dev \
        libfftw3-dev \
        make \
        patch \
        python3.5-dev \
        python3-pip \
        python3-setuptools \
    && apt-get clean -y \
    && apt-get autoclean -y \
    && rm -rf /var/lib/apt/lists/* /tmp/*

COPY root /
ENTRYPOINT [ "/usr/bin/earthworm" ]
CMD [ ]
