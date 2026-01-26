class Moji < Formula
  desc "Terminal art toolkit - kaomoji, ASCII banners, filters, QR codes, and more"
  homepage "https://github.com/ddmoney420/moji"
  version "1.0.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_darwin_arm64.tar.gz"
      sha256 "ab104fdd2d7fc7bfc33464c6d0651a297e5c61e69172697caeb283effa20d209"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_darwin_amd64.tar.gz"
      sha256 "eb1f6cc94388879fb5942d4dbbd69f595f9ba208a324c501c1b700c8c026ce70"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_linux_arm64.tar.gz"
      sha256 "d2e6048fe3f28df9dfa31215da1d900498e1f25aac0521b9efc5c7f54f039c67"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_linux_amd64.tar.gz"
      sha256 "3d61635dd8a3b018b1c7d9a61a8a192d888975b4063ef0e348b73e72a537173e"
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
