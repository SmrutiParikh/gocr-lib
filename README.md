This Project uses Golang to play with images using OpenCV, ImageMagick and Tesseract.

### TL;DR: Running Examples

- **For Text Extraction**:
    ```bash
    make run PLAIN_TEXT_EXTRACTION samples/documents/Eric_BROOKS-Resume.jpg eng
    make run PLAIN_TEXT_EXTRACTION samples/documents/japanese.png jpn
    ```

- **For Text Extraction using HOCR**:
    ```bash
    make run HOCR_TEXT_EXTRACTION samples/documents/Eric_BROOKS-Resume.jpg eng
    make run HOCR_TEXT_EXTRACTION samples/documents/japanese jpn
    ```

- **For Image Object Detection**:
    ```bash
    make run IMAGE_OBJECT_DETECTION samples/images/traffic.jpg eng
    ```

- **For Video Object Detection**:
    ```bash
    make run VIDEO_OBJECT_DETECTION samples/videos/marathon.mp4 eng
    ```

# Project Setup Guide

This project utilizes **OpenCV**, **ImageMagick**, and **Tesseract** for image processing and text extraction. Below are the steps to install the dependencies, set up the Go project, integrate OpenCV for image processing, and use Tesseract for text extraction.

## 1. Install Go 1.22.x

#### **For Linux (Ubuntu/Debian)**

1. **Download the Go 1.22.x tarball:**
   ```bash
   wget https://golang.org/dl/go1.22.2.linux-amd64.tar.gz
   ```

2. **Extract the tarball:**
   ```bash
   sudo tar -C /usr/local -xvzf go1.22.2.linux-amd64.tar.gz
   ```

3. **Set Go environment variables:**
   Add the following lines to your `.bashrc` or `.zshrc` file:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   export GOPATH=$HOME/go
   ```

4. **Apply the changes:**
   ```bash
   source ~/.bashrc  # or source ~/.zshrc
   ```

5. **Verify the installation:**
   ```bash
   go version
   ```
   This should display `go version go1.22.x linux/amd64`.

#### **For macOS**

1. **Download Go 1.22.x package for macOS:**
   - Visit [Go Downloads](https://golang.org/dl/) and download the `go1.22.x.darwin-amd64.pkg` file.

2. **Install Go**:
   - Double-click the `.pkg` file to start the installation process and follow the on-screen instructions.

3. **Verify the installation:**
   ```bash
   go version
   ```
   This should display `go version go1.22.x darwin/amd64`.

#### **For Windows**

1. **Download the Go 1.22.x installer:**
   - Visit [Go Downloads](https://golang.org/dl/) and download the `go1.22.x.windows-amd64.msi` installer.

2. **Run the installer**:
   - Double-click the `.msi` file to start the installation process and follow the on-screen instructions.

3. **Verify the installation**:
   Open Command Prompt and run:
   ```bash
   go version
   ```
   This should display `go version go1.22.x windows/amd64`.


## 2. Install Dependencies

### OpenCV (4.10.0)
OpenCV is used for image processing tasks like object detection and manipulation. Since we are using `gocv` version `0.39.0`, which is compatible with OpenCV 4.10.0, here are the steps to manually install OpenCV.

Ref: https://github.com/hybridgroup/gocv/releases

#### **Steps to Install OpenCV (4.10.0)**

**Step1: Download OpenCV and opencv_contrib from GitHub**:
1. Download OpenCV (4.10.0):
    - Visit the [OpenCV GitHub page](https://github.com/opencv/opencv/releases/tag/4.10.0).
    - Download the appropriate release for your system (e.g., `.tar.gz` for Linux/Mac or `.zip` for Windows).

    Download opencv_contrib:

2. Go to [the opencv_contrib GitHub repository](https://github.com/opencv/opencv_contrib).
    Download the opencv_contrib repository(4.x).

**Step2: Extract the files**:
1. For Linux/Mac:
    - Extract OpenCV:
    ```bash
    tar -xvzf opencv-4.10.0.tar.gz
    ```
    - Extract opencv_contrib:
    ```bash
    tar -xvzf opencv_contrib.tar.gz
    ```

2. For Windows, use a tool like 7-Zip to extract the `.zip` files.

**Step3: Install Dependencies**:
- On Linux, run the following commands to install required libraries:
    ```bash
    sudo apt-get install -y build-essential cmake git pkg-config libgtk-3-dev libcanberra-gtk* libxvidcore-dev libx264-dev libjpeg-dev libpng-dev libtiff-dev libjasper-dev libopenexr-dev libdcmtk-dev libpdc-dev libfftw3-dev
    sudo apt-get install -y libeigen3-dev
    ```

**Step4: Build OpenCV with opencv_contrib**:
1. Create a build directory:
      ```bash
      cd opencv-4.10.0
      mkdir build
      cd build
      ```
2. Configure the build with cmake:
    - Run cmake to configure OpenCV with opencv_contrib:
        ```bash
        cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_PREFIX=/usr/local -D OPENCV_EXTRA_MODULES_PATH=../opencv_contrib/modules ..
        ```
3. Compile and Install OpenCV
    - Run make to compile OpenCV and opencv_contrib:
        ```bash
        make -j$(nproc)
        sudo make install
        ```
4. Verify the Installation
    - Verify that OpenCV and opencv_contrib were installed correctly by running
        ```bash
        pkg-config --modversion opencv4
        python -c "import cv2.aruco; print('ArUco module available')"
        python -c "import cv2; print(cv2.__version__)"
        ```
        OR
        ```bash
        pkg-config --modversion opencv4
        python3 -c "import cv2.aruco; print('ArUco module available')"
        python3 -c "import cv2; print(cv2.__version__)"
        ```
        Should return 
        ```
        4.10.0
        ArUco module available
        4.10.0
        ```

**Step5: Link OpenCV with Go**:
  Make sure gocv is linked to the OpenCV installation:

    ```bash
    go get -u gocv.io/x/gocv@0.39.0
    ```

### ImageMagick (7.x)
ImageMagick is used for image manipulation, such as resizing, rotating, or converting images.

#### **Steps to Install ImageMagick (7.x)**

1. **Download ImageMagick 7 from GitHub**:
    - Visit the [ImageMagick GitHub Releases page](https://github.com/ImageMagick/ImageMagick/releases).
    - Since gocv supports ImageMagic 7, download the appropriate stable version.

2. **Install Dependencies**:
    - For Linux:
      ```bash
      sudo apt-get install -y build-essential
      sudo apt-get install -y libjpeg-dev libpng-dev libtiff-dev libfreetype6-dev
      ```

3. **Build and Install**:
    - Extract the ImageMagick archive, and navigate into the extracted directory.
      ```bash
      tar -xvzf ImageMagick-7.x.x.tar.gz
      cd ImageMagick-7.x.x
      ```
    - Run the following commands to configure and install ImageMagick:
      ```bash
      ./configure
      make
      sudo make install
      sudo ldconfig /usr/local/lib
      ```

4. **Verify Installation**:
    - Run:
      ```bash
      convert --version
      ```
    - This should display the installed version of ImageMagick.

### Tesseract
Tesseract is an OCR (Optical Character Recognition) engine used to extract text from images.

#### **Install Tesseract**

1. **Install Tesseract**:
    - On Linux:
      ```bash
      sudo apt-get install tesseract-ocr
      ```

2. **Verify Installation**:
    - Run:
      ```bash
      tesseract --version
      ```
    - This should display the installed version of Tesseract.

---

## 3. Set Up Go Project

Since you already have a `go.mod` file with the required dependencies, you can set up your Go project as follows:

1. **Clone the repository (if applicable):**

    ```bash
    git clone <repository_url>
    cd <project_directory>
    ```
2. **Ensure the dependencies are installed**: To fetch all the necessary dependencies for your Go project:

    ```bash
    go mod tidy
    ```


## 4. Running the Project

1. **Build the project**:
    ```bash
    make build
    ```

2. **Run the project**:
    For text extraction:
    ```bash
    make run TEXT_EXTRACTION samples/bill.jpg
    ```
    For object detection:
    ```bash
    go run main.go OBJECT_DETECTION samples/kangaroo_horse.png
    ```

3. **Clean the build**:
    ```bash
    make clean
    ```

4. **Test the project**:
    ```bash
    make test
    ```




