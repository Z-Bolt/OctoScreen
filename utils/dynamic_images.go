package utils

import (
	"github.com/gotk3/gotk3/gtk"
)


func CreateFilamentSpoolImage(color string) (*gtk.Image, error) {
	svg :=
		"<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n" +
		"<svg version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" x=\"0px\" y=\"0px\" width=\"62.227px\" height=\"62.889px\" viewBox=\"0 0 62.227 62.889\" enable-background=\"new 0 0 62.227 62.889\" xml:space=\"preserve\">\n" +
		"	<g id=\"Filament_layer\">\n" +
		"		<g id=\"Filament\">\n" +
		"			<path fill=\"" + color + "\" d=\"M31.188,2.401C15.068,2.401,2,15.469,2,31.589s13.067,29.188,29.188,29.188S60.376,47.71,60.376,31.589 S47.308,2.401,31.188,2.401z M31.188,45.215c-7.524,0-13.625-6.101-13.625-13.625c0-7.525,6.101-13.625,13.625-13.625 s13.625,6.1,13.625,13.625S38.713,45.215,31.188,45.215z\"/>\n" +
		"		</g>\n" +
		"	</g>\n" +
		"	<g id=\"Spool_layer\">\n" +
		"		<g id=\"Hex_wall\">\n" +
		"			<g id=\"Hex_fill\">\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.569\" y1=\"8.493\" x2=\"9.85\" y2=\"13.103\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,22.335 3.889,21.173 5.897,16.993 8.569,13.119 8.572,13.119 9.881,13.119 12.543,17.726 9.881,22.335\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,31.551 1.9,26.941 4.56,22.337 9.881,22.337 12.543,26.941 9.881,31.551\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,40.767 1.9,36.159 4.56,31.552 9.881,31.552 12.543,36.159 9.881,40.767\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,40.769 3.889,41.931 5.897,46.109 8.569,49.984 8.572,49.984 9.881,49.984 12.543,45.377 9.881,40.769\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"9.834\"  y1=\"49.978\" x2=\"12.553\" y2=\"54.588\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"8.503\"  x2=\"17.875\" y2=\"8.503\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"17.722\" x2=\"17.875\" y2=\"17.722\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"26.942\" x2=\"17.875\" y2=\"26.942\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"36.162\" x2=\"17.875\" y2=\"36.162\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"45.383\" x2=\"17.875\" y2=\"45.383\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"54.602\" x2=\"17.875\" y2=\"54.602\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"20.56,13.117 17.9,8.508 20.56,3.902 25.881,3.902 28.543,8.508 25.881,13.117\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"21.174,22.335 20.56,22.335 17.9,17.726 20.56,13.119 25.881,13.119 28.543,17.726 28.238,18.253\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"18.19,27.443 17.9,26.941 20.56,22.337 21.154,22.337\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"21.113,40.767 20.56,40.767 17.9,36.159 18.219,35.608\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"28.262,44.891 28.543,45.375 25.881,49.984 20.56,49.984 17.9,45.375 20.56,40.769 21.154,40.769\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"20.56,59.2 17.9,54.594 20.56,49.985 25.881,49.985 28.543,54.594 25.881,59.2\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"8.503\"  x2=\"33.875\" y2=\"8.503\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"17.722\" x2=\"33.875\" y2=\"17.722\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"45.383\" x2=\"33.875\" y2=\"45.383\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"54.602\" x2=\"33.875\" y2=\"54.602\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"36.561,13.117 33.901,8.508 36.561,3.902 41.881,3.902 44.543,8.508 41.881,13.117\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"34.211,18.263 33.901,17.726 36.561,13.119 41.881,13.119 44.543,17.726 41.881,22.335 41.254,22.335\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"41.222,22.337 41.881,22.337 44.543,26.941 44.217,27.506\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"44.206,35.575 44.543,36.159 41.881,40.767 41.238,40.767\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"41.254,40.769 41.881,40.769 44.543,45.375 41.881,49.984 36.561,49.984 33.901,45.375 34.189,44.874\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"36.561,59.2 33.901,54.594 36.561,49.985 41.881,49.985 44.543,54.594 41.881,59.2\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"8.503\"  x2=\"49.875\" y2=\"8.503\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"17.722\" x2=\"49.875\" y2=\"17.722\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"26.942\" x2=\"49.875\" y2=\"26.942\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"36.162\" x2=\"49.875\" y2=\"36.162\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"45.383\" x2=\"49.875\" y2=\"45.383\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"54.602\" x2=\"49.875\" y2=\"54.602\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"49.834\" y1=\"8.478\"  x2=\"52.569\" y2=\"13.118\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"57.884,22.335 58.555,21.173 56.547,16.994 53.875,13.119 53.873,13.119 52.563,13.119 49.901,17.726 52.563,22.335\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"52.561,31.551 49.901,26.941 52.561,22.337 57.881,22.337 60.543,26.941 57.881,31.551\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"52.561,40.767 49.901,36.159 52.561,31.552 57.881,31.552 60.543,36.159 57.881,40.767\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"57.884,40.769 58.555,41.931 56.547,46.109 53.875,49.984 53.873,49.984 52.563,49.984 49.901,45.377 52.563,40.769\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"52.553\" y1=\"49.963\" x2=\"49.834\" y2=\"54.588\"/>\n" +
		"			</g>\n" +
		"			<circle id=\"Hex_inner_ring\" fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.562\"  stroke-miterlimit=\"10\" cx=\"31.188\" cy=\"31.59\"  r=\"13.625\"/>\n" +
		"		</g>\n" +
		"		<g id=\"Rings\">\n" +
		"			<circle id=\"Outter_ring\"    fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"3.0024\" stroke-miterlimit=\"10\" cx=\"31.188\" cy=\"31.589\" r=\"29.188\"/>\n" +
		"			<circle id=\"Inner_ring\"     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"3.0024\" stroke-miterlimit=\"10\" cx=\"31.188\" cy=\"31.59\"  r=\"13.625\"/>\n" +
		"		</g>\n" +
		"	</g>\n" +
		"</svg>"
		
	return ImageNewFromSvg(svg)
}

func CreateFilamentSpoolWithCheckMarkImage(color string) (*gtk.Image, error) {
	svg :=
		"<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n" +
		"<svg version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" x=\"0px\" y=\"0px\" width=\"62.227px\" height=\"62.889px\" viewBox=\"0 0 62.227 62.889\" enable-background=\"new 0 0 62.227 62.889\" xml:space=\"preserve\">\n" +
		"	<g id=\"Filament_layer\">\n" +
		"		<g id=\"Filament\">\n" +
		"			<path fill=\"" + color + "\" d=\"M31.188,2.401C15.068,2.401,2,15.469,2,31.589s13.066,29.188,29.188,29.188S60.376,47.71,60.376,31.589 S47.308,2.401,31.188,2.401z M31.188,45.215c-7.523,0-13.625-6.101-13.625-13.625c0-7.525,6.102-13.625,13.625-13.625 c7.524,0,13.625,6.1,13.625,13.625S38.713,45.215,31.188,45.215z\"/>\n" +
		"		</g>\n" +
		"	</g>\n" +
		"	<g id=\"Spool_layer\">\n" +
		"		<g id=\"Hex_wall\">\n" +
		"			<g id=\"Hex_fill\">\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.569\" y1=\"8.493\" x2=\"9.85\" y2=\"13.103\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,22.335 3.889,21.173 5.897,16.993 8.569,13.119 8.572,13.119 9.881,13.119 12.543,17.726 9.881,22.335\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,31.551 1.9,26.941 4.56,22.337 9.881,22.337 12.543,26.941 9.881,31.551\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,40.767 1.9,36.159 4.56,31.552 9.881,31.552 12.543,36.159 9.881,40.767\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"4.56,40.769 3.889,41.931 5.897,46.109 8.569,49.984 8.572,49.984 9.881,49.984 12.543,45.377 9.881,40.769\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"9.834\"  y1=\"49.978\" x2=\"12.553\" y2=\"54.588\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"8.503\"  x2=\"17.875\" y2=\"8.503\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"17.722\" x2=\"17.875\" y2=\"17.722\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"26.942\" x2=\"17.875\" y2=\"26.942\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"36.162\" x2=\"17.875\" y2=\"36.162\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"45.383\" x2=\"17.875\" y2=\"45.383\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"12.531\" y1=\"54.602\" x2=\"17.875\" y2=\"54.602\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"20.56,13.117 17.9,8.508 20.56,3.902 25.881,3.902 28.543,8.508 25.881,13.117\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"21.174,22.335 20.56,22.335 17.9,17.726 20.56,13.119 25.881,13.119 28.543,17.726 28.238,18.253\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"18.19,27.443 17.9,26.941 20.56,22.337 21.154,22.337\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"21.113,40.767 20.56,40.767 17.9,36.159 18.219,35.608\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"28.262,44.891 28.543,45.375 25.881,49.984 20.56,49.984 17.9,45.375 20.56,40.769 21.154,40.769\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"20.56,59.2 17.9,54.594 20.56,49.985 25.881,49.985 28.543,54.594 25.881,59.2\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"8.503\"  x2=\"33.875\" y2=\"8.503\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"17.722\" x2=\"33.875\" y2=\"17.722\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"45.383\" x2=\"33.875\" y2=\"45.383\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"28.531\" y1=\"54.602\" x2=\"33.875\" y2=\"54.602\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"36.561,13.117 33.901,8.508 36.561,3.902 41.881,3.902 44.543,8.508 41.881,13.117\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"34.211,18.263 33.901,17.726 36.561,13.119 41.881,13.119 44.543,17.726 41.881,22.335 41.254,22.335\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"41.222,22.337 41.881,22.337 44.543,26.941 44.217,27.506\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"44.206,35.575 44.543,36.159 41.881,40.767 41.238,40.767\"/>\n" +
		"				<polyline fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"41.254,40.769 41.881,40.769 44.543,45.375 41.881,49.984 36.561,49.984 33.901,45.375 34.189,44.874\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"36.561,59.2 33.901,54.594 36.561,49.985 41.881,49.985 44.543,54.594 41.881,59.2\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"8.503\"  x2=\"49.875\" y2=\"8.503\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"17.722\" x2=\"49.875\" y2=\"17.722\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"26.942\" x2=\"49.875\" y2=\"26.942\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"36.162\" x2=\"49.875\" y2=\"36.162\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"45.383\" x2=\"49.875\" y2=\"45.383\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"44.532\" y1=\"54.602\" x2=\"49.875\" y2=\"54.602\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"49.834\" y1=\"8.478\"  x2=\"52.569\" y2=\"13.118\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"57.884,22.335 58.555,21.173 56.547,16.994 53.875,13.119 53.873,13.119 52.563,13.119 49.901,17.726 52.563,22.335\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"52.561,31.551 49.901,26.941 52.561,22.337 57.881,22.337 60.543,26.941 57.881,31.551\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"52.561,40.767 49.901,36.159 52.561,31.552 57.881,31.552 60.543,36.159 57.881,40.767\"/>\n" +
		"				<polygon  fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" points=\"57.884,40.769 58.555,41.931 56.547,46.109 53.875,49.984 53.873,49.984 52.563,49.984 49.901,45.377 52.563,40.769\"/>\n" +
		"				<line     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.5616\" stroke-miterlimit=\"10\" x1=\"52.553\" y1=\"49.963\" x2=\"49.834\" y2=\"54.588\"/>\n" +
		"			</g>\n" +
		"			<circle id=\"Hex_inner_ring\" fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"0.562\"  stroke-miterlimit=\"10\" cx=\"31.188\" cy=\"31.59\"  r=\"13.625\"/>\n" +
		"		</g>\n" +
		"		<g id=\"Rings\">\n" +
		"			<circle id=\"Outter_ring\"    fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"3.0024\" stroke-miterlimit=\"10\" cx=\"31.188\" cy=\"31.589\" r=\"29.188\"/>\n" +
		"			<circle id=\"Inner_ring\"     fill=\"none\" stroke=\"#FFFFFF\" stroke-width=\"3.0024\" stroke-miterlimit=\"10\" cx=\"31.188\" cy=\"31.59\"  r=\"13.625\"/>\n" +
		"		</g>\n" +
		"	</g>\n" +
		"	<g id=\"Check_mark_layer\">\n" +
		"		<circle id=\"Circle\" fill=\"#FFFFFF\" cx=\"48\" cy=\"14.594\" r=\"13.594\"/>\n" +
		"		<path id=\"Check_mark\" d=\"M54.88,8.738l-8.991,8.991l-3.292-3.291l-2.306,2.302l5.598,5.556l11.301-11.32L54.88,8.738z\"/>\n" +
		"	</g>\n" +
		"</svg>"
		
	return ImageNewFromSvg(svg)
}

func Create_X_Image(color string) (*gtk.Image, error) {
	svg :=
		"<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n" +
		""

	return ImageNewFromSvg(svg)
}
