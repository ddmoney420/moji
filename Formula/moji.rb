class Moji < Formula
  desc "Terminal art toolkit - kaomoji, ASCII banners, filters, QR codes, and more"
  homepage "https://github.com/ddmoney420/moji"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_darwin_arm64.tar.gz"
      sha256 "581d1b0d0ca40a7a2dd65458f3f79e82db7737e4ef6094bca7dd912763ddd7bd"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_darwin_amd64.tar.gz"
      sha256 "af8a2d25dcbe81646416f85d4b6a52f96aa447e8d895c343c8cf8648e74f06e9"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_linux_arm64.tar.gz"
      sha256 "2c2f8d9456be0889fab62863043b3164a9be7d1294e5c0f6b4c6e5c22dbe92a6"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_linux_amd64.tar.gz"
      sha256 "aa804595ab57d659a6aa4cd917c0c0c04938df92be4d5a2039a85b148f2b6485"
    end
  end

  def install
    bin.install "moji"
  end

  test do
    assert_match "moji", shell_output("#{bin}/moji --version")
    assert_match "shrug", shell_output("#{bin}/moji list | head -20")
  end
end
