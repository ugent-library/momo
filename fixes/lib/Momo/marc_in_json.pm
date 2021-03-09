package Momo::marc_in_json;

use Catmandu::Sane;
use Moo;
use Catmandu::Fix::Has;

has key => (fix_arg => 1);

sub fix {
    my ($self, $data) = @_;

    if (my $marc = delete $data->{$self->key}) {
        my $mij = $data->{$self->key} = {};
        for my $field (@$marc) {
            my ($tag, $ind1, $ind2, @subfields) = @$field;
 
            if ($tag eq 'LDR') {
                shift @subfields;
                $mij->{leader} = join "", @subfields;
            }
            elsif ($tag eq 'FMT' || substr($tag, 0, 2) eq '00') {
                shift @subfields;
                push @{$mij->{fields} ||= []},
                    {$tag => join "", @subfields};
            }
            else {
                my @sf;
                my $start = !defined($subfields[0])
                    || $subfields[0] eq '_' ? 2 : 0;
                for (my $i = $start; $i < @subfields; $i += 2) {
                    push @sf, {$subfields[$i] => $subfields[$i + 1]};
                }
                push @{$mij->{fields} ||= []},
                    {$tag => {subfields => \@sf, ind1 => $ind1, ind2 => $ind2}};
            }
        }
    }
 
    $data;
}

1;
