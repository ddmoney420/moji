class Moji < Formula
  desc "Terminal art toolkit - kaomoji, ASCII banners, filters, QR codes, and more"
  homepage "https://github.com/ddmoney420/moji"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_#{version}_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_#{version}_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_#{version}_linux_arm64.tar.gz"
      sha256 "PLACEHOLDER"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_#{version}_linux_amd64.tar.gz"
      sha256 "PLACEHOLDER"
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
