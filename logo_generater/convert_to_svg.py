import os
import subprocess
import sys
import argparse
from pathlib import Path


def check_dependencies():
    """检查系统是否已安装必要的转换依赖"""
    deps = []
    
    # 检查 pdf2svg
    try:
        subprocess.run(['pdf2svg', '--help'], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        deps.append('pdf2svg')
    except (subprocess.CalledProcessError, FileNotFoundError):
        pass # 可选依赖
    
    # 检查 pdftocairo (Poppler)
    try:
        subprocess.run(['pdftocairo', '-v'], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        deps.append('pdftocairo')
    except (subprocess.CalledProcessError, FileNotFoundError):
        pass # 可选依赖
    
    return deps


def pdf_to_svg(pdf_path, svg_path=None):
    """
    尝试多种方法将 PDF 转换为 SVG
    """
    if svg_path is None:
        svg_path = Path(pdf_path).with_suffix('.svg')
    else:
        # 如果提供了输出路径，确保是 Path 对象
        svg_path = Path(svg_path)
    
    # 确保输入是 Path 对象
    pdf_path = Path(pdf_path)

    print(f"正在将 PDF 转换为 SVG: {pdf_path} -> {svg_path}")

    if not pdf_path.exists():
        print(f"❌ 错误: 输入的 PDF 文件不存在: {pdf_path}")
        return None
    
    # 方法1: 尝试使用 pdf2svg 工具
    if can_use_tool('pdf2svg'):
        try:
            result = subprocess.run([
                'pdf2svg', 
                str(pdf_path), 
                str(svg_path)
            ], check=True, capture_output=True, text=True)
            
            print("✓ 使用 pdf2svg 转换成功！")
            return svg_path
            
        except (subprocess.CalledProcessError, FileNotFoundError):
            print("⚠ pdf2svg 转换失败，尝试其他方法...")
    
    # 方法2: 尝试使用 poppler 的 pdftocairo 工具
    if can_use_tool('pdftocairo'):
        try:
            # pdftocairo -svg input.pdf output_prefix (它会自动添加 .svg)
            # 所以我们需要剥离后缀作为前缀
            output_prefix = str(svg_path.with_suffix(''))
            
            result = subprocess.run([
                'pdftocairo', 
                '-svg',
                str(pdf_path),
                output_prefix
            ], check=True, capture_output=True, text=True)
            
            # 检查生成的文件
            # pdftocairo 总是附加 .svg
            actual_svg_path = Path(output_prefix + '.svg')
            
            if actual_svg_path.exists():
                # 如果实际生成的和我们想要的名字（例如大小写差异或路径差异）不一样，重命名
                if actual_svg_path.resolve() != svg_path.resolve():
                    if svg_path.exists():
                         svg_path.unlink() # 覆盖
                    actual_svg_path.rename(svg_path)

                print("✓ 使用 pdftocairo 转换成功！")
                return svg_path
            else:
                raise Exception("预期生成的SVG文件未找到")
                
        except (subprocess.CalledProcessError, FileNotFoundError, Exception) as e:
            print(f"⚠ pdftocairo 转换失败: {e}，尝试其他方法...")
    
    # 方法3: 使用 PyMuPDF (fitz) 作为最后手段
    return convert_with_pymupdf(pdf_path, svg_path)


def can_use_tool(tool_name):
    """检查是否可以使用某个工具"""
    try:
        if tool_name == 'pdf2svg':
            subprocess.run(['pdf2svg', '--help'], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        elif tool_name == 'pdftocairo':
            subprocess.run(['pdftocairo', '-v'], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False


def convert_with_pymupdf(pdf_path, svg_path):
    """
    使用 PyMuPDF 将 PDF 转换为 SVG
    """
    print("正在尝试使用 PyMuPDF 进行转换...")
    
    try:
        import fitz  # PyMuPDF
    except ImportError:
        print("PyMuPDF 未安装，正在尝试安装...")
        try:
            subprocess.check_call([sys.executable, '-m', 'pip', 'install', 'pymupdf'])
            import fitz
        except Exception as e:
            print(f"✗ 无法安装 PyMuPDF: {e}")
            return None
    
    try:
        # 打开 PDF 文件
        pdf_document = fitz.open(pdf_path)
        # 获取第一页
        page = pdf_document[0]
        # 转换为 SVG
        svg_data = page.get_svg_image()
        
        # 保存 SVG 文件 (覆盖模式)
        with open(svg_path, 'w', encoding='utf-8') as f:
            f.write(svg_data)
        
        print(f"✓ 使用 PyMuPDF 转换成功: {svg_path}")
        return svg_path
        
    except Exception as e:
        print(f"✗ PyMuPDF 转换失败: {e}")
        return None


def main():
    parser = argparse.ArgumentParser(
        description='将 PDF 转换为 SVG 格式',
        formatter_class=argparse.RawTextHelpFormatter,
        epilog="""使用示例:
  python convert_to_svg.py                    # 默认转换 logo.pdf -> logo.svg
  python convert_to_svg.py myicon.pdf         # 转换指定文件
  python convert_to_svg.py -o icon.svg        # 指定输出文件名"""
    )
    parser.add_argument('input_pdf', nargs='?', help='输入的 PDF 文件路径 (默认: logo.pdf)')
    parser.add_argument('-o', '--output', help='输出的 SVG 文件路径')
    
    args = parser.parse_args()
    
    # 默认文件
    input_pdf = args.input_pdf or 'logo.pdf'
    
    # 执行前的依赖检查日志可选
    # check_dependencies()
    
    # 直接执行转换
    svg_result = pdf_to_svg(input_pdf, args.output)
    
    if svg_result:
        print(f"\n🎉 成功生成 SVG 文件: {svg_result}")
        return 0
    else:
        print("\n❌ 转换失败")
        return 1


if __name__ == '__main__':
    exit(main())